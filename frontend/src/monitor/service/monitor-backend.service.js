(function () {
    'use strict';
    angular.module('app')
        .factory('monitorBackend', monitorBackend);
    /* @ngInject */
    function monitorBackend($resource) {
        return {
            monitor: monitor
        };

        function monitor(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/query', {
                cluster: paramObj.cluster,
                user: paramObj.user,
                metric: paramObj.metric,
                app: paramObj.app,
                task: paramObj.task,
                start: paramObj.start,
                end: paramObj.end,
                step: paramObj.step,
                expr: paramObj.expr
            });
        }
    }
})();
