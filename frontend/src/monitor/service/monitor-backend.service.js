(function () {
    'use strict';
    angular.module('app')
        .factory('monitorBackend', monitorBackend);
    /* @ngInject */
    function monitorBackend($resource) {
        return {
            listApp: listApp,
            monitor: monitor
        };

        function listApp() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/applications');
        }

        function monitor(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor', {
                appid: paramObj.appid,
                metric: paramObj.metric,
                type: paramObj.type,
                from: paramObj.from,
                to: paramObj.to,
                step: paramObj.step
            });
        }
    }
})();
