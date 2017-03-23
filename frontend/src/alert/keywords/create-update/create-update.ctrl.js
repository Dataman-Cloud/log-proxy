(function () {
    'use strict';
    angular.module('app')
        .controller('CreateAlertKeywordCtrl', CreateAlertKeywordCtrl);
    /* @ngInject */
    function CreateAlertKeywordCtrl(alertBackend, $state, target, alert, Notification) {
        var alert = alert.data || {};

        var self = this;

        self.target = target;
        self.form = {
            app: alert.app,
            keyword: alert.keyword || '',
            source: alert.source || ''
        };

        self.create = create;
        self.update = update;

        activate();

        function activate() {

        }

        function create() {
            alertBackend.alerts().save(self.form, function (data) {
                Notification.success('创建成功');
                $state.go('home.alertkeyword', null, {reload: true})
            });
        }

        function update() {
            self.form.id = alert.id;
            alertBackend.alerts().update(self.form, function (data) {
                Notification.success('更新成功');
                $state.go('home.alertkeyword', null, {reload: true})
            });
        }
    }
})();
