(function () {
    'use strict';
    angular.module('app')
        .controller('EventMonitorCtrl', EventMonitorCtrl);
    /* @ngInject */
    function EventMonitorCtrl(eventBackend) {
        var self = this;

        self.monitors = [];

        activate();

        function activate() {
            listEvents()
        }

        function listEvents() {
            eventBackend.monitor().get(function (data) {
                self.monitors = data.data.results;
            })
        }
    }
})();