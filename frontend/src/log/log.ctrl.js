(function () {
    'use strict';
    angular.module('app')
        .controller('LogCtrl', LogCtrl);
    /* @ngInject */
    function LogCtrl(logBackend, moment) {
        var self = this;
        self.curTimestamp = moment().unix() * 1000;
        self.fromTimestamp = moment().subtract(1, 'hours').unix() * 1000;

        activate();

        function activate() {
            logBackend.listApp({from: self.fromTimestamp, to: self.curTimestamp}).get(function (data) {
                console.log(data)
            })
        }
    }
})();