var constants = require('../constants');

var Sequelize = require('sequelize');
var sqldb = require('./mysql/DBModel');
var modelDB = sqldb.CrowdFunding;

exp = {
    getById: function (id) {
        return modelDB.findOne({
            order: [
                ['id', 'DESC']],
            where: {
                id: id,
            }
        });
    },
    getByHId: function (hid) {
        return modelDB.findOne({
            order: [
                ['id', 'DESC']],
            where: {
                hid: hid,
            }
        });
    },
    updateActiveHid: function (tx, id, hid) {
        return modelDB.update(
            {
                hid: hid,
                status: constants.CROWD_STATUS_APPROVED,
                date_actived: new Date()
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateActiveFailed: function (tx, id) {
        return modelDB.update(
            {
                status: constants.CROWD_STATUS_APPROVED_FAILED,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateShakedInfo: function (tx, id, balance, qty) {
        return modelDB.update(
            {
                balance: balance,
                shaked_num: modelDB.sequelize.literal('shake_num + ' + qty),
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateFunded: function (tx, id) {
        return modelDB.update(
            {
                status: constants.CROWD_STATUS_FUNDED,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateFundedFailed: function (tx, id) {
        return modelDB.update(
            {
                status: constants.PROJECT_STATUS_FUNDED_FAILED,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateDelivered: function (tx, id) {
        return modelDB.update(
            {
                status: constants.PROJECT_STATUS_DELIVERED,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateDeliveredFailed: function (tx, id) {
        return modelDB.update(
            {
                status: constants.PROJECT_STATUS_DELIVERED_FAILED,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateCanceled: function (tx, id) {
        return modelDB.update(
            {
                status: constants.PROJECT_STATUS_CANCELED,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateCanceledFailed: function (tx, id) {
        return modelDB.update(
            {
                status: constants.PROJECT_STATUS_CANCELED_FAILED,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    getExpiredActivedProjects: function () {
        return modelDB.findAll({
            order: [
                ['id', 'DESC']],
            where: {
                end_date: {
                    [Sequelize.Op.lt]: Sequelize.literal('now()')
                },
                status: constants.PROJECT_STATUS_APPROVED,
            }
        });
    },
    updateStatus: function (id, status) {
        return modelDB.update(
            {
                status: status,
            },
            {
                where:
                    {
                        id: id
                    }
            }
        );
    },
}

module.exports = exp;