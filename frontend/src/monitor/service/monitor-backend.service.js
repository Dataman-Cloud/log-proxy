(function () {
    'use strict';
    angular.module('app')
        .factory('monitorBackend', monitorBackend);
    /* @ngInject */
    function monitorBackend($resource) {
        return {
            monitor: monitor,
            clusters: clusters,
            apps: apps,
            tasks: tasks,
            metrics: metrics
        };

        function monitor(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/query', {
                cluster: paramObj.cluster,
                metric: paramObj.metric,
                app: paramObj.app,
                task: paramObj.task,
                start: paramObj.start,
                end: paramObj.end,
                step: paramObj.step,
                expr: paramObj.expr
            });
        }

        function clusters(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/clusters', {
                start: paramObj.start,
                end: paramObj.end,
                step: paramObj.step
            });
        }

        function apps(cluster, data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/clusters/:cluster/apps', {
                cluster: cluster,
                start: paramObj.start,
                end: paramObj.end,
                step: paramObj.step
            });
        }

        function tasks(cluster, app, data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/clusters/:cluster/apps/:app/tasks', {
                cluster: cluster,
                app: app,
                start: paramObj.start,
                end: paramObj.end,
                step: paramObj.step
            });
        }
        
        function metrics() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/query/metrics');
        }
    }
})();
