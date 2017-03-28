/**
 * Created by lixiaoyan on 2017/3/27.
 */
/**
 * Created by lixiaoyan on 2017/3/24.
 */
(function () {
    'use strict';
    angular.module('app')
        .controller('AlertRulesCtrl', AlertRulesCtrl);
    /* @ngInject */

    function AlertRulesCtrl(monitorAlertBackend, $stateParams, events, timing, Notification, $state, $scope) {
        var self = this;

        self.timePeriod = 60;
        self.count = 0;
        var appListReloadInterval = 10000;
        self.events = events.data.events;
        self.form = {
            cluster: $stateParams.cluster || '',
            app: $stateParams.app || '',
            start: moment().subtract(self.timePeriod, 'minutes').unix(),
            end: moment().unix(),
            page: 1,
            size: 100,
            ack: false
        };


        self.searchAlert = searchAlert;
        self.activateAck = activateAck;

        timing.start($scope, reloadSearch, appListReloadInterval);

        activate();

        function activate() {

        }

        function search() {
            monitorAlertBackend.events(self.form).get(function (data) {
                self.events = data.data.events;
                self.count = data.data.count;
            })
        }

        function reloadSearch() {
            return monitorAlertBackend.events(self.form).get()
                .$promise.then(function (data) {
                    self.events = data.data.events;
                    self.count = data.data.count;
                })
        }

        function searchAlert() {
            if (!self.timePeriod) {
                search({page: 1, size: 100});
            } else {
                self.form.size = 100;
                self.form.page = 1;

                self.form.end = moment().unix();
                self.form.start = moment().subtract(self.timePeriod, 'minutes').unix();

                search(self.form);
            }

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