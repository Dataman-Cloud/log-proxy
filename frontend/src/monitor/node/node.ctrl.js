(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorNodeCtrl', MonitorNodeCtrl);
    /* @ngInject */
    function MonitorNodeCtrl(nodes) {
        var self = this;
        self.nodes = nodes.data.nodes;

        activate();

        function activate() {

        }
    }
})();