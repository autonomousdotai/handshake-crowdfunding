var Sequelize = require('sequelize');

var sqldb = require('./mysql/DBModel');
var modelDB = sqldb.CrowdFundingShaked;
var constants = require('../constants');

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
    checkUserBacked: function (crowd_funding_id, user_id, id) {
        return modelDB.findOne({
            order: [
                ['id', 'DESC']],
            where: {
                crowd_funding_id: crowd_funding_id,
                user_id: user_id,
                id: {
                    [Sequelize.Op.ne]: id
                },
                status: {
                    [Sequelize.Op.gt]: constants.CROWD_ORDER_STATUS_NEW
                },
            }
        });
    },
    getFirstBacked: function (neuron_project_id, customer_id) {
        return modelDB.findOne({
            order: [
                ['id', 'ASC']],
            where: {
                neuron_project_id: neuron_project_id,
                customer_id: customer_id,
                status: {
                    [Sequelize.Op.in]: [constants.PROJECT_ORDER_STATUS_APPROVED, constants.PROJECT_ORDER_STATUS_UNSHAKED_PROCESS]
                },
            }
        });
    },
    updateActived: function (tx, id, address) {
        address = address.toLowerCase()
        return modelDB.update(
            {
                status: constants.CROWD_ORDER_STATUS_APPROVED,
                address: address,
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
                status: constants.CROWD_ORDER_STATUS_APPROVED_FAILED,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateUserUnshaked: function (tx, user_id, crowd_funding_id) {
        return modelDB.update(
            {
                status: constants.CROWD_ORDER_STATUS_UNSHAKED,
            },
            {
                where:
                    {
                        user_id: user_id,
                        crowd_funding_id: crowd_funding_id,
                        status: constants.CROWD_ORDER_STATUS_UNSHAKED_PROCESS,
                    }
            }, {transaction: tx}
        );
    },
    updateUserCanceled: function (tx, user_id, crowd_funding_id) {
        return modelDB.update(
            {
                status: constants.CROWD_ORDER_STATUS_CANCELED,
            },
            {
                where:
                    {
                        customer_id: user_id,
                        neuron_project_id: crowd_funding_id,
                        status: constants.CROWD_ORDER_STATUS_CANCELED_PROCESS,
                    }
            }, {transaction: tx}
        );
    },
    updateUserRefunded: function (tx, user_id, crowd_funding_id) {
        return modelDB.update(
            {
                status: constants.CROWD_ORDER_STATUS_REFUNDED,
            },
            {
                where:
                    {
                        customer_id: user_id,
                        neuron_project_id: crowd_funding_id,
                        status: constants.CROWD_ORDER_STATUS_REFUNDED_PROCESS,
                    }
            }, {transaction: tx}
        );
    },
    updateUserUnshakeFailed: function (tx, user_id, crowd_funding_id) {
        return modelDB.update(
            {
                status: constants.CROWD_ORDER_STATUS_UNSHAKED_FAILED,
            },
            {
                where:
                    {
                        user_id: user_id,
                        crowd_funding_id: crowd_funding_id,
                        status: constants.CROWD_ORDER_STATUS_UNSHAKED_PROCESS,
                    }
            }, {transaction: tx}
        );
    },
    updateUserCancelFailed: function (tx, user_id, crowd_funding_id) {
        return modelDB.update(
            {
                status: constants.CROWD_ORDER_STATUS_CANCELED_FAILED,
            },
            {
                where:
                    {
                        user_id: user_id,
                        crowd_funding_id: crowd_funding_id,
                        status: constants.CROWD_ORDER_STATUS_CANCELED_PROCESS,
                    }
            }, {transaction: tx}
        );
    },
    updateUserRefundFailed: function (tx, user_id, crowd_funding_id) {
        return modelDB.update(
            {
                status: constants.CROWD_ORDER_STATUS_REFUNDED_FAILED,
            },
            {
                where:
                    {
                        user_id: user_id,
                        crowd_funding_id: crowd_funding_id,
                        status: constants.CROWD_ORDER_STATUS_REFUNDED_PROCESS,
                    }
            }, {transaction: tx}
        );
    },
}

module.exports = exp;