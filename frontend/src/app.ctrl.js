(function () {
    'use strict';
    angular.module('app')
        .controller('RootCtrl', RootCtrl);

    /* @ngInject */
    function RootCtrl($window, $state) {
        var self = this;

        self.goBack = goBack;

        activate();

        function activate() {

        }

        function goBack(state) {
            if (state) {
                $state.go(state);
            } else {
                $window.history.length > 2 ? $window.history.back() : $state.go('app');
            }
        }

    }
})();
