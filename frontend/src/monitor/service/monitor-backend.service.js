(function () {
    'use strict';
    angular.module('app')
        .factory('monitorBackend', monitorBackend);
    /* @ngInject */
    function monitorBackend($resource) {
        return {
            listApp: listApp,
            listInstance: listInstance,
            monitor: monitor,
            listNode: listNode
        };

        function listApp() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/applications');
        }

        function listInstance(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/application', {
                appid: paramObj.appid
            });
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

        function listNode() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/nodes');
        }
    }
})();
