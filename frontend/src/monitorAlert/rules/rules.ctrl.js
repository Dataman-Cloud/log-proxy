(function () {
    'use strict';
    angular.module('app')
        .controller('AlertRulesCtrl', AlertRulesCtrl);
    /* @ngInject */

    function AlertRulesCtrl(monitorAlertBackend, $stateParams, timing, Notification, $state, $scope) {
        var self = this;
        var appListReloadInterval = 10000;

        self.timePeriod = ($stateParams.start && $stateParams.end) ?
            ($stateParams.end - $stateParams.start) / 60 : '';

        self.count = 0;

        self.events = [];
        self.form = {
            cluster: $stateParams.cluster || '',
            app: $stateParams.app || '',
            start: $stateParams.start || '',
            end: $stateParams.end || '',
            page: 1,
            size: 100,
            ack: false
        };


        self.searchAlert = searchAlert;
        self.activateAck = activateAck;

        timing.start($scope, reloadSearch, appListReloadInterval);

        activate();

        function activate() {
            search()
        }

        function search() {
            return monitorAlertBackend.events(self.form).get().$promise.then(function (data) {
                self.events = data.data.events;
                self.count = data.data.count;

                return data
            })

        }

        function reloadSearch() {
            if (self.timePeriod) {
                self.form = {
                    cluster: $stateParams.cluster || '',
                    app: $stateParams.app || '',
                    start: moment().subtract(self.timePeriod, 'minutes').unix(),
                    end: moment().unix(),
                    page: 1,
                    size: 100,
                    ack: false
                };
            } else {
                self.form = {
                    cluster: $stateParams.cluster || '',
                    app: $stateParams.app || '',
                    page: 1,
                    size: 100,
                    ack: false
                };
            }
            return monitorAlertBackend.events(self.form).get()
                .$promise.then(function (data) {
                    self.events = data.data.events;
                    self.count = data.data.count;
                });
        }

        function searchAlert() {
            if (self.timePeriod) {
                self.form = {
                    cluster: $stateParams.cluster || '',
                    app: $stateParams.app || '',
                    start: moment().subtract(self.timePeriod, 'minutes').unix(),
                    end: moment().unix(),
                    page: 1,
                    size: 100,
                    ack: false
                };
            } else {
                self.form = {
                    cluster: $stateParams.cluster || '',
                    app: $stateParams.app || '',
                    page: 1,
                    size: 100,
                    ack: false
                };
            }
            $state.go('home.alertRules', self.form, {inherit: false})
        }

        function activateAck(event) {
            var obj = {
                action: 'ack',
                cluster: event.cluster,
                app: event.app
            };

            monitorAlertBackend.eventsAck({id: event.ID}).update(obj, function (data) {
                Notification.success('操作成功');
                $state.reload();
            });

        }

    }
})();