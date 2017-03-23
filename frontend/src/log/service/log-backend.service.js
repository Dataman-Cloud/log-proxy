(function () {
    'use strict';
    angular.module('app')
        .factory('logBackend', logBackend);
    /* @ngInject */
    function logBackend($resource) {
        return {
            searchLogs: searchLogs,
            logContext: logContext,
            clusters: clusters,
            apps: apps,
            tasks: tasks,
            paths: paths
        };

        function searchLogs(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/index', {
                cluster: paramObj.cluster,
                app: paramObj.app,
                task: paramObj.task,
                path: paramObj.path,
                keyword: paramObj.keyword,
                from: paramObj.from,
                to: paramObj.to,
                page: paramObj.page,
                size: paramObj.size
            });
        }

        function logContext(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/context', {
                cluster: paramObj.cluster,
                app: paramObj.app,
                task: paramObj.task,
                path: paramObj.path,
                offset: paramObj.offset,
                page: paramObj.page,
                size: paramObj.size
            });
        }
        
        function clusters(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/clusters', {
                from: paramObj.from,
                to: paramObj.to
            });
        }

        function apps(cluster, data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/clusters/:cluster/apps', {
                cluster: cluster,
                from: paramObj.from,
                to: paramObj.to
            });
        }
        
        function tasks(cluster, app, data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/clusters/:cluster/apps/:app/tasks', {
                cluster: cluster,
                app: app,
                from: paramObj.from,
                to: paramObj.to
            });
        }

        function paths(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/paths/:app', {
                cluster: paramObj.cluster,
                app: paramObj.app,
                task: paramObj.task,
                from: paramObj.from,
                to: paramObj.to
            });
        }
    }
})();
