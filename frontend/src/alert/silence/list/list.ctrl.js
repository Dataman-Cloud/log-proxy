/**
 * Created by lixiaoyan on 2016/12/1.
 */
(function () {
    'use strict';
    angular.module('app')
        .controller('AlertSilencesCtrl', AlertSilencesCtrl);
    /* @ngInject */
    function AlertSilencesCtrl(alertBackend, alertcurd) {
        var curDate = new Date();
        var self = this;

        self.silences = [];
        self.silenceDetailSet = {};

        self.deleteSilence = deleteSilence;
        self.checkState = checkState;

        activate();

        function activate() {
            listSilent()
        }

        function listSilent() {
            alertBackend.silences().get(function (data) {
                self.silences = data.data
            })
        }

        function deleteSilence(id) {
            alertcurd.deleteSilence(id)
        }

        function checkState(startsAt, endsAt) {
            if (new Date(startsAt) < curDate && new Date(endsAt) > curDate) {
                return 'Active'
            } else if (new Date(startsAt) > curDate) {
                return 'Pending'
            } else {
                return 'Elapsed'
            }
        }

    }
})();
