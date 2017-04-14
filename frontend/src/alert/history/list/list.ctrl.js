(function () {
    'use strict';
    angular.module('app')
        .controller('AlertHistoriesCtrl', AlertHistoriesCtrl);
    /* @ngInject */
    function AlertHistoriesCtrl(alertBackend, $stateParams, moment, Notification, $state) {
        var self = this;

        self.timePeriod = 60;

        self.count = 0;
        self.histories = [];

        self.clusters = [];
        self.apps = [];

        self.form = {
            cluster: $stateParams.cluster || '',
            app: $stateParams.app || '',
            keyword: $stateParams.keyword,
            start: moment().subtract(self.timePeriod, 'minutes').unix(),
            end: moment().unix(),
            page: 1,
            size: 100
        };

        //md-table parameter
        self.options = {
            rowSelection: true,
            decapitate: false,
            boundaryLinks: false,
            limitSelect: true,
            pageSelect: true
        };

        self.defaultLimit = 100;
        self.limitOptions = [100, 200, 500];
        self.query = {
            limit: 100,
            page: 1
        };

        self.loadClusters = loadClusters;
        self.loadApps = loadApps;
        self.onPaginate = onPaginate;
        self.searchHistory = searchHistory;
        self.isAck = isAck;
        self.infoContext = infoContext;

        activate();

        function activate() {
            loadClusters();
            loadApps();
            listHistory();
        }

        function loadClusters() {
            checkTimeRange();

            alertBackend.clusters({start: self.startTime, end: self.endTime}).get(function (data) {
                self.clusters = data.data;
            })
        }

        function loadApps() {
            if (self.form.cluster) {
                alertBackend.apps({
                    cluster: self.form.cluster,
                    start: self.startTime,
                    end: self.endTime
                }).get(function (data) {
                    self.apps = data.data;
                })
            }
        }

        function fetchHistory(data) {
            alertBackend.histories(data).get(function (data) {
                if (data.data) {
                    self.histories = data.data.events;
                    self.count = data.data.count;

                }
            })
        }

        function searchHistory() {
            if (!self.timePeriod) {
                fetchHistory({page: 1, size: 100});
            } else {
                self.form.size = 100;
                self.form.page = 1;

                self.form.end = moment().unix();
                self.form.start = moment().subtract(self.timePeriod, 'minutes').unix();

                fetchHistory(self.form);
            }

        }

        function listHistory() {
            if ($stateParams.cluster && $stateParams.app) {
                self.form.size = 100;
                self.form.page = 1;

                fetchHistory(self.form)
            }
        }

        function onPaginate(page, limit) {
            self.form.size = limit;
            self.form.page = page;

            fetchHistory(self.form)
        }

        function checkTimeRange() {
            if (self.timePeriod) {
                self.startTime = moment().subtract(self.timePeriod, 'minutes').unix();
                self.endTime = moment().unix();
            } else {
                self.startTime = null;
                self.endTime = null;
            }
        }

        function isAck(history) {
            var obj = {
                action: 'ack'
            };
            alertBackend.historyAck({id: history.id}).pupdate(obj, function (data) {
                history.ack = true;
                Notification.success('操作成功');
            });

        }

        function infoContext(history) {
            $state.go('home.logbase.logcontext', {
                app: history.appid,
                task: history.taskid,
                source: history.path,
                offset: history.offset,
                cluster: history.clusterid,
                start: moment(history.logtime+'').subtract(1,'h').unix()*1000,  // 1 小时前
                end: moment(history.logtime+'').add(1,'h').unix()*1000  // 1 小时后
            })
        }
    }
})();
