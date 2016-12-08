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
        self.alerts = [];

        self.deleteKeyword = deleteKeyword;

        activate();

        function activate() {
            listKeyword()
        }

        function listKeyword() {
            alertBackend.alerts().get(function (data) {
                self.keywords = data.data.results
            })
        }

        function deleteKeyword(id) {
            alertcurd.deleteKeyword(id)
        }

    }
})();
