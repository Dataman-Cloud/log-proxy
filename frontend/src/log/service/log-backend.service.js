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
            sources: sources
        };

        function searchLogs(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/clusters/:cluster/apps/:app/index', {
                cluster: paramObj.cluster,
                app: paramObj.app,
                task: paramObj.task,
                source: paramObj.source,
                keyword: paramObj.keyword,
                from: paramObj.from,
                to: paramObj.to,
                page: paramObj.page,
                size: paramObj.size,
                conj: paramObj.conj
            });
        }

        function logContext(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/clusters/:cluster/apps/:app/context', {
                cluster: paramObj.cluster,
                app: paramObj.app,
                task: paramObj.task,
                source: paramObj.source,
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

        function sources(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/clusters/:cluster/apps/:app/sources', {
                cluster: paramObj.cluster,
                app: paramObj.app,
                task: paramObj.task,
                from: paramObj.from,
                to: paramObj.to
            });
        }
    }
})();
