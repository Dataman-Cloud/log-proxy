(function () {
    'use strict';
    angular.module('app')
        .controller('CreateAlertCtrl', CreateAlertCtrl);
    /* @ngInject */
    function CreateAlertCtrl(alertBackend, $state, target, alert) {
        var alert = alert.data || {};

        var self = this;

        self.target = target;
        self.form = {
            appid: alert.appid,
            keyword: alert.keyword || '',
            path: alert.path || ''
        };

        self.create = create;
        self.update = update;

        activate();

        function activate() {

        }

        function create() {
            alertBackend.alerts().save(self.form, function (data) {
                $state.go('home.alertkeyword', null, {reload: true})
            });
        }

        function update() {
            self.form.id = alert.id;
            alertBackend.alerts().update(self.form, function (data) {
                $state.go('home.alertkeyword', null, {reload: true})
            });
        }
    }
})();
