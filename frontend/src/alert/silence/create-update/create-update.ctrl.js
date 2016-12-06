(function () {
    'use strict';
    angular.module('app')
        .controller('CreateAlertSilenceCtrl', CreateAlertSilenceCtrl);
    /* @ngInject */
    function CreateAlertSilenceCtrl(alertBackend, $state, $stateParams, target, silence, Notification) {
        var self = this;

        self.target = target;
        self.target === 'create' ? self.form = {
            createdBy: '',
            comment: '',
            endsAt: '',
            startsAt: '',
            matchers: [{
                name: '',
                value: ''
            }]
        } : self.form = formatSilence(silence.data);

        self.create = create;
        self.update = update;
        self.addMatcher = addMatcher;
        self.deleteMatcher = deleteMatcher;

        activate();

        function activate() {

        }

        function formatSilence(silence) {
            silence.startsAt = new Date(silence.startsAt);
            silence.endsAt = new Date(silence.endsAt);
            return silence
        }

        function addMatcher() {
            var matcher = {
                name: '',
                value: ''
            };

            self.form.matchers.push(matcher);
        }

        function deleteMatcher(index) {
            self.form.matchers.splice(index, 1);
        }

        function create() {
            alertBackend.silences().save(self.form, function (data) {
                Notification.success('创建成功');
                $state.go('home.alertSilences', null, {reload: true})
            })
        }

        function update() {
            if ($stateParams.from === 'silence') {
                var silenceId = self.form.id;
                delete self.form.id;
                alertBackend.silence(silenceId).update(self.form, function (data) {
                    Notification.success('更新成功');
                    $state.go('home.alertSilences', null, {reload: true})
                })
            }
        }
    }
})();
