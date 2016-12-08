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
                appid: paramObj.appid,
                metric: paramObj.metric,
                start: paramObj.start,
                end : paramObj.end,
                step: paramObj.step,
                taskid: paramObj.taskid,
                expr: paramObj.expr
            });
        }

    }
})();
