/**
 * Created by lixiaoyan on 2016/12/1.
 */
(function () {
    'use strict';
    angular.module('app')
        .controller('AlertKeywordsCtrl', AlertKeywordsCtrl);
    /* @ngInject */
    function AlertKeywordsCtrl(alertBackend, alertcurd) {
        var self = this;

        self.ALERT_STATUS = {
            Enabled: '激活',
            Disabled: '未激活'
        };

        self.count = 0;
        self.keywords = [];

        //md-table parameter
        self.options = {
            rowSelection: true,
            decapitate: false,
            boundaryLinks: false,
            limitSelect: true,
            pageSelect: true
        };

        self.defaultLimit = 20;
        self.limitOptions = [20, 50, 100];
        self.query = {
            limit: self.defaultLimit,
            page: 1
        };

        self.deleteKeyword = deleteKeyword;
        self.activateAlarm = activateAlarm;
        self.onPaginate = onPaginate;

        activate();

        function activate() {
            listKeyword()
        }

        function listKeyword() {
            alertBackend.alerts({size: self.defaultLimit, page: 1}).get(function (data) {
                self.keywords = data.data.rules;
                self.count = data.data.count;
            })
        }

        function deleteKeyword(id) {
            alertcurd.deleteKeyword(id)
        }

        function activateAlarm(data) {
            alertcurd.activateAlarm(data)
        }

        function onPaginate(page, limit) {
            alertBackend.alerts({size: limit, page: page}).get(function (data) {
                self.keywords = data.data.rules;
                self.count = data.data.count;
            })
        }

    }
})();
