/**
 * Created by lixiaoyan on 2016/12/1.
 */
(function () {
    'use strict';
    angular.module('app')
        .controller('AlertListCtrl', AlertListCtrl);
    /* @ngInject */
    function AlertListCtrl(alertBackend, alertcurd) {
        var self = this;
        self.alerts = [];

        self.deleteAlert = deleteAlert;

        activate();

        function activate() {
            listAlert()
        }
        
        function listAlert() {
            alertBackend.alerts().get(function (data) {
                self.alerts = data.data
            })
        }

        function deleteAlert(id) {
            alertcurd.deleteAlert(id)
        }

    }
})();
