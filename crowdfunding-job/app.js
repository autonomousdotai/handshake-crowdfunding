var cron = require('node-cron');

var ethEventDAO = require('./models/eth_event_dao');
var crowdFundingDAO = require('./models/crowd_funding_dao');
var crowdFundingShakeDAO = require('./models/crowd_funding_shake_dao');
var ethTxDAO = require('./models/eth_tx_dao');


var sqldb = require('./models/mysql/DBModel');
var modelDB = sqldb.db;

var config = require('./configs/index')
var constants = require('./constants');


var Web3 = require('web3');
web3 = new Web3(new Web3.providers.HttpProvider(config.blockchainNetwork));

var CrowdsaleContract = require('./contracts/CrowdsaleHandshake.json');
var crowdsaleContractAddress = config.crowdsaleContractAddress
var crowdsaleContractEventNames = ['__init', '__shake', '__unshake', '__cancel', '__refund', '__stop', '__withdraw'];
var crowdsaleContractIns = new web3.eth.Contract(CrowdsaleContract.abi, crowdsaleContractAddress);

console.log('Events by blockchainNetwork: ' + config.blockchainNetwork);
console.log('Events by crowdsaleContractAddress: ' + crowdsaleContractAddress);

function parseOffchain(offchain) {
    let values = offchain.replace(/\u0000/g, '').split("_")
    console.log(values)
    if (values.length >= 2) {
        return [values[0].trim(), values[1].trim()];
    } else {
        return null;
    }
}

async function processEventObj(contractAddress, eventName, eventObj) {
    let tx = await modelDB.transaction();
    try {
        console.log("processEventObj", contractAddress, eventName, eventObj);

        let tx_hash = eventObj.transactionHash.toLowerCase()
        let txr = await web3.eth.getTransactionReceipt(tx_hash);

        await ethEventDAO.create(tx, contractAddress, eventName, JSON.stringify(eventObj), eventObj.blockNumber, eventObj.logIndex);

        switch (contractAddress) {
            case crowdsaleContractAddress: {
                switch (eventName) {
                    case '__init': {
                        console.log("__init hid = " + eventObj.returnValues.hid);
                        console.log("__init offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || offchain == undefined) {
                            console.log("__init missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__init offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_CROWD) {
                            let offchainId = parseInt(offchains[1]);
                            let crowdFunding = await crowdFundingDAO.getById(offchainId);
                            if (crowdFunding == null) {
                                console.log("__init crowdFundingDAO.getById NULL", offchainId);
                                break;
                            }
                            console.log("__init crowdFundingDAO.getById OK", offchainId);
                            if (crowdFunding.id > 0 && crowdFunding.status == constants.CROWD_STATUS_NEW) {
                                await crowdFundingDAO.updateActiveHid(tx, crowdFunding.id, hid);
                                console.log("__init crowdFundingDAO.updateActiveHid OK", crowdFunding.id, hid);
                            }
                            let tns = await ethTxDAO.getByHash(tx_hash);
                            if (tns == null) {
                                tns = await ethTxDAO.create(crowdFunding.user_id, tx_hash, 'crowd_init', crowdFunding.id);
                            }
                        }
                    }
                        break;
                    case '__shake': {
                        console.log("__shake hid = " + eventObj.returnValues.hid);
                        console.log("__shake state = " + eventObj.returnValues.state);
                        console.log("__shake balance = " + eventObj.returnValues.balance);
                        console.log("__shake offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const state = eventObj.returnValues.state;
                        const balance = eventObj.returnValues.balance;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || state == undefined || balance == undefined || offchain == undefined) {
                            console.log("__shake missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__shake offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_SHAKED) {
                            let crowdFundingId = parseInt(offchains[1]);
                            let crowdFundingShaked = await crowdFundingShakeDAO.getById(crowdFundingId);
                            if (crowdFundingShaked == null) {
                                console.log("__shake crowdFundingShakedDAO.getById NULL", crowdFundingId);
                                break;
                            }
                            console.log("__shake crowdFundingShakedDAO.getById OK", crowdFundingId);
                            if (crowdFundingShaked.id > 0 && crowdFundingShaked.status == constants.CROWD_ORDER_STATUS_NEW) {
                                let address = '';
                                if (txr != null) {
                                    address = txr.from
                                }
                                await crowdFundingShakeDAO.updateActived(tx, crowdFundingShaked.id, address)
                                console.log("__shake crowdFundingShakedDAO.updateActived OK", crowdFundingShaked.id, address)

                                let crowdFunding = await crowdFundingDAO.getById(crowdFundingShaked.crowd_funding_id);
                                if (crowdFunding != null && crowdFunding.status > 0) {
                                    console.log("__shake crowdFundingDAO.getById OK", crowdFundingShaked.crowd_funding_id);

                                    const balanceEth = Web3.utils.fromWei(balance, 'ether');
                                    let balanceEthU = balanceEth
                                    let crowdFundingShakedCheck = await crowdFundingShakeDAO.checkUserBacked(crowdFunding.id, crowdFundingShaked.user_id, crowdFundingShaked.id)
                                    let qtyU = 0;
                                    if (crowdFundingShakedCheck == null) {
                                        qtyU = 1;
                                        console.log("__shake crowdFundingShakedDAO.checkUserBacked NULL", crowdFunding.id, crowdFundingShaked.user_id, crowdFundingShaked.id)
                                    } else {
                                        console.log("__shake crowdFundingShakedDAO.checkUserBacked OK", crowdFunding.id, crowdFundingShaked.user_id, crowdFundingShaked.id)
                                    }

                                    await crowdFundingDAO.updateShakedInfo(tx, crowdFunding.id, balanceEthU, qtyU);
                                    console.log("__shake crowdFundingDAO.updateShakedInfo OK", crowdFunding.id, balanceEthU, qtyU);

                                    if (state == constants.CROWSALE_STATE_SHAKED) {
                                        if (crowdFunding.id > 0 && crowdFunding.status == constants.CROWD_STATUS_FAILED && crowdFunding.crowd_date < new Date()) {
                                            await crowdFundingDAO.updateFunded(tx, crowdFunding.id);
                                            console.log("__shake crowdFundingDAO.updateFunded OK", crowdFunding.id);
                                        }
                                    }
                                }
                                let ethTx = await ethTxDAO.getByHash(tx_hash);
                                if (ethTx == null) {
                                    ethTx = await ethTxDAO.create(crowdFundingShaked.user_id, tx_hash, 'crowd_shake', crowdFundingShaked.id, txr.from, txr.to, new Date(), txr.blockNumber, 0, 0, JSON.stringify(txr));
                                }
                            }
                        }
                    }
                        break;
                    case '__unshake': {
                        console.log("__unshake hid = " + eventObj.returnValues.hid);
                        console.log("__unshake state = " + eventObj.returnValues.state);
                        console.log("__unshake balance = " + eventObj.returnValues.balance);
                        console.log("__unshake offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const state = eventObj.returnValues.state;
                        const balance = eventObj.returnValues.balance;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || state == undefined || balance == undefined || offchain == undefined) {
                            console.log("__unshake missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__unshake offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_USER) {
                            let userId = parseInt(offchains[1]);

                            let crowdFunding = await crowdFundingDAO.getByHId(hid);
                            if (crowdFunding != null && crowdFunding.status > 0) {
                                console.log("__unshake crowdFundingDAO.getByHId OK", hid);

                                const balanceEth = Web3.utils.fromWei(balance, 'ether');
                                await crowdFundingDAO.updateShakedInfo(tx, crowdFunding.id, balanceEth, -1);
                                console.log("__unshake crowdFundingDAO.updateShakedInfo OK", crowdFunding.id, balanceEth, -1);

                                await crowdFundingShakeDAO.updateUserUnshaked(tx, userId, crowdFunding.id);
                                console.log("__unshake crowdFundingShakedDAO.updateUserUnshaked OK", userId, crowdFunding.id);

                                let ethTx = await ethTxDAO.getByHash(tx_hash);
                                if (ethTx == null) {
                                    ethTx = await ethTxDAO.create(userId, tx_hash, 'crowd_unshake', crowdFunding.id, txr.from, txr.to, new Date(), txr.blockNumber, 0, 0, JSON.stringify(txr));
                                }
                            }
                        }
                    }
                        break;
                    case '__cancel': {
                        console.log("__cancel hid = " + eventObj.returnValues.hid);
                        console.log("__cancel state = " + eventObj.returnValues.state);
                        console.log("__cancel offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const state = eventObj.returnValues.state;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || state == undefined || offchain == undefined) {
                            console.log("__cancel missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__cancel offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_USER) {
                            let userId = parseInt(offchains[1]);
                            let crowdFunding = await crowdFundingDAO.getByHId(hid);
                            if (crowdFunding != null) {
                                console.log("__cancel crowdFundingDAO.getByHId OK", hid);

                                await crowdFundingShakeDAO.updateUserCanceled(tx, userId, crowdFunding.id);
                                console.log("__cancel crowdFundingShakedDAO.updateUserCanceled OK", userId, crowdFunding.id);

                                if (crowdFunding.id > 0 && state == constants.CROWSALE_STATE_CANCELED && crowdFunding.status != constants.CROWD_STATUS_CANCELED) {
                                    await crowdFundingDAO.updateCanceled(tx, crowdFunding.id)
                                    console.log("__cancel crowdFundingDAO.updateCanceled OK", crowdFunding.id);
                                }

                                let ethTx = await ethTxDAO.getByHash(tx_hash);
                                if (ethTx == null) {
                                    ethTx = await ethTxDAO.create(userId, tx_hash, 'crowd_cancel', crowdFunding.id, txr.from, txr.to, new Date(), txr.blockNumber, 0, 0, JSON.stringify(txr));
                                }
                            }
                        }
                    }
                        break;
                    case '__refund': {
                        console.log("__refund hid = " + eventObj.returnValues.hid);
                        console.log("__refund state = " + eventObj.returnValues.state);
                        console.log("__refund offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const state = eventObj.returnValues.state;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || state == undefined || offchain == undefined) {
                            console.log("__refund missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__refund offchains", offchains);
                        let offchainType = offchains[0];

                        if (offchainType == constants.OFFCHAIN_TYPE_USER) {
                            let userId = parseInt(offchains[1]);

                            let crowdFunding = await crowdFundingDAO.getByHId(hid);
                            if (crowdFunding != null && crowdFunding.status > 0) {
                                console.log("__refund crowdFundingDAO.getByHId OK", hid);

                                await crowdFundingShakeDAO.updateUserRefunded(tx, userId, crowdFunding.id);
                                console.log("__refund crowdFundingShakedDAO.updateUserRefunded OK", userId, crowdFunding.id);

                                let ethTx = await ethTxDAO.getByHash(tx_hash);
                                if (ethTx == null) {
                                    ethTx = await ethTxDAO.create(userId, tx_hash, 'crowd_refund', crowdFunding.id, txr.from, txr.to, new Date(), txr.blockNumber, 0, 0, JSON.stringify(txr));
                                }
                            }
                        }
                    }
                        break;
                    // case '__stop': {
                    //     console.log("__stop hid = " + eventObj.returnValues.hid);
                    //     console.log("__refund state = " + eventObj.returnValues.state);
                    //     console.log("__stop offchain = " + eventObj.returnValues.offchain);
                    //     const hid = eventObj.returnValues.hid;
                    //     const state = eventObj.returnValues.state;
                    //     const offchain = eventObj.returnValues.offchain;
                    //     if (hid == undefined || state == undefined || offchain == undefined) {
                    //         console.log("__stop missing parameters");
                    //         break;
                    //     }
                    //     let offchainStr = Web3.utils.toAscii(offchain);
                    //     let offchains = parseOffchain(offchainStr);
                    //     console.log("__stop offchains", offchains);
                    //     let offchainType = offchains[0];
                    //     if (offchainType == constants.OFFCHAIN_TYPE_PROJECT) {
                    //         let offchainId = parseInt(offchains[1]);
                    //         let project = await crowdFundingDAO.getById(offchainId);
                    //         if (project != null && project.status > 0) {
                    //             if (project.id > 0 && state == constants.CROWSALE_STATE_CANCELED && project.status != constants.PROJECT_STATUS_CANCELED) {
                    //                 await crowdFundingDAO.updateCanceled(tx, project.id)
                    //                 console.log("__stop crowdFundingDAO.updateCanceled OK", project.id);
                    //             }
                    //             let tns = await ethTxDAO.getByHash(tx_hash);
                    //             if (tns == null) {
                    //                 tns = await ethTxDAO.create(project.customer_id, tx_hash, 'stopProject', project.id);
                    //             }
                    //             sendEmailFundCanceledCreator(project.id);
                    //             sendEmailFundCanceledBacker(project.id);
                    //         }
                    //     }
                    // }
                    //     break;
                    case '__withdraw': {
                        console.log("__withdraw hid = " + eventObj.returnValues.hid);
                        console.log("__withdraw amount = " + eventObj.returnValues.amount);
                        console.log("__withdraw offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const amount = eventObj.returnValues.amount;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || amount == undefined || offchain == undefined) {
                            console.log("__withdraw missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__withdraw offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_CROWD) {
                            let crowdFundingId = parseInt(offchains[1]);
                            let crowdFunding = await crowdFundingDAO.getById(crowdFundingId);
                            if (crowdFunding != null && crowdFunding.status > 0) {
                                let ethTx = await ethTxDAO.getByHash(tx_hash);
                                if (ethTx == null) {
                                    ethTx = await ethTxDAO.create(crowdFunding.user_id, tx_hash, 'crowd_withdraw', crowdFunding.id, txr.from, txr.to, new Date(), txr.blockNumber, 0, 0, JSON.stringify(txr));
                                }
                            }
                        }
                    }
                        break;
                }
            }
                break;
        }
        tx.commit();
    } catch (err) {
        console.log('processEventObj', err);
        tx.rollback();
    }
}

function asyncGetPastEvents(contract, contractAddress, eventName, fromBlock) {
    return new Promise(function (resolve, reject) {
        contract.getPastEvents(eventName, {
            filter: {_from: contractAddress},
            fromBlock: fromBlock,
            toBlock: 'latest'

        }, function (error, events) {
            console.log(eventName + " getPastEvents OK")
            if (error != null) {
                reject(error);
            } else {
                resolve(events);
            }
        });
    })
}

async function asyncScanEventLog(contract, contractAddress, eventName) {
    let lastEventLog = await ethEventDAO.getLastLogByName(contractAddress, eventName);
    var fromBlock = 0;
    if (lastEventLog != null) {
        fromBlock = lastEventLog.block + 1;
    }
    console.log(eventName + " fromBlock = " + fromBlock);
    let events = await asyncGetPastEvents(contract, contractAddress, eventName, fromBlock);
    for (var i = 0; i < events.length; i++) {
        const eventObj = events[i];
        console.log(eventObj);
        let checkEventLog = await ethEventDAO.getByBlock(contractAddress, eventObj.blockNumber, eventObj.logIndex);
        if (checkEventLog == null) {
            await processEventObj(contractAddress, eventName, eventObj);
        }
    }

}

async function processTx(id, user_id, hash, ref_type, ref_id, date_created) {
    let tx = await modelDB.transaction();
    try {
        let txr = null;
        try {
            txr = await web3.eth.getTransactionReceipt(hash);
        } catch (err) {
            console.log('error', err)
            txr = null;
        }
        let is_failed = false
        if (txr == null) {
            let now = new Date()
            if (now - date_created > 24 * 60 * 60 * 1000) {
                is_failed = true
            } else {
                await ethTxDAO.updateStatus(tx, hash, 0);
                console.log('txr is null', hash);
                tx.commit();
                return
            }
        } else {
            console.log('txr is ok', txr);
            let txrJson = JSON.stringify(txr);
            await ethTxDAO.updateInfo(tx, id, txr.from, txr.to, new Date(), txr.blockNumber, 0, 0, txrJson);
            is_failed = (txr.status == '1' || txr.status == '0x1') ? false : true;
        }
        if (is_failed) {
            if (txr != null) {
                await ethTxDAO.updateStatus(tx, hash, 2);
            } else {
                await ethTxDAO.updateStatus(tx, hash, 3);
            }
        } else {
            await ethTxDAO.updateStatus(tx, hash, 1);
        }
        switch (ref_type) {
            case 'crowd_init':{
                if (is_failed) {
                    let crowdFunding = await crowdFundingDAO.getById(ref_id);
                    if (crowdFunding == null) {
                        console.log(ref_type + ' crowdFundingDAO.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' crowdFundingDAO.getById OK', ref_id);
                    await crowdFundingDAO.updateActiveFailed(tx, crowdFunding.id);
                    console.log(ref_type + ' crowdFundingDAO.updateNewFailed OK', crowdFunding.id);
                }
            }
                break;
            case 'crowd_shake':{
                if (is_failed) {
                    let crowdFundingShaked = await crowdFundingShakeDAO.getById(ref_id);
                    if (crowdFundingShaked == null) {
                        console.log(ref_type + ' crowdFundingShakedDAO.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' crowdFundingShakedDAO.getById OK', ref_id);
                    await crowdFundingShakeDAO.updateActiveFailed(tx, crowdFundingShaked.id);
                    console.log(ref_type + ' crowdFundingShakedDAO.updateNewFailed OK', crowdFundingShaked.id);
                }
            }
                break;
            case 'crowd_unshake':{
                if (is_failed) {
                    let crowdFunding = await crowdFundingDAO.getById(ref_id);
                    if (crowdFunding == null) {
                        console.log(ref_type + ' crowdFundingDAO.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' crowdFundingDAO.getById OK', ref_id);
                    if (crowdFunding.id > 0) {
                        await crowdFundingShakeDAO.updateUserUnshakeFailed(tx, user_id, crowdFunding.id);
                        console.log("__cancel crowdFundingShakedDAO.updateUserUnshakeFailed OK", user_id, crowdFunding.id);
                    }
                }
            }
                break;
            case 'crowd_cancel':{
                if (is_failed) {
                    let crowdFunding = await crowdFundingDAO.getById(ref_id);
                    if (crowdFunding == null) {
                        console.log(ref_type + ' crowdFundingDAO.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' crowdFundingDAO.getById OK', ref_id);
                    if (crowdFunding.id > 0) {
                        await crowdFundingShakeDAO.updateUserCancelFailed(tx, user_id, crowdFunding.id);
                        console.log("cancelProject crowdFundingShakedDAO.updateUserCancelFailed OK", user_id, crowdFunding.id);
                    }
                }
            }
                break;
            case 'crowd_refund':{
                if (is_failed) {
                    let crowdFunding = await crowdFundingDAO.getById(ref_id);
                    if (crowdFunding == null) {
                        console.log(ref_type + ' crowdFundingDAO.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' crowdFundingDAO.getById OK', ref_id);
                    if (crowdFunding.id > 0) {
                        await crowdFundingShakeDAO.updateUserRefundFailed(tx, user_id, crowdFunding.id);
                        console.log("refundProject crowdFundingShakedDAO.updateUserRefundFailed OK", user_id, crowdFunding.id);
                    }
                }
            }
                break;
            case 'crow_stop':{
                if (is_failed) {
                    let crowdFunding = await crowdFundingDAO.getById(ref_id);
                    if (crowdFunding == null) {
                        console.log(ref_type + ' crowdFundingDAO.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' crowdFundingDAO.getById OK', ref_id);
                    if (crowdFunding.id > 0) {
                        await crowdFundingDAO.updateCanceledFailed(tx, crowdFunding.id);
                        console.log(ref_type + ' crowdFundingDAO.updateCanceledFailed OK', crowdFunding.id);
                    }
                }
            }
                break;
        }
        tx.commit();
    } catch (err) {
        console.log('error', err)
        tx.rollback();
    }
}

async function cronJob() {
    console.log('running a task every minute at ' + new Date());
    console.log('process ether tx');
    let results = await ethTxDAO.getListUnTx();
    for (var i = 0; i < results.length; i++) {
        var result = results[i];
        await processTx(result.id, result.user_id, result.hash, result.ref_type, result.ref_id, result.date_created);
    }
    console.log('process ether events');
    if (crowdsaleContractAddress != '') {
        for (var i = 0; i < crowdsaleContractEventNames.length; i++) {
            var eventName = crowdsaleContractEventNames[i];
            await asyncScanEventLog(crowdsaleContractIns, crowdsaleContractAddress, eventName);
        }
    }
}

cronJob();

cron.schedule('* * * * *', async function () {
    await cronJob();
});

