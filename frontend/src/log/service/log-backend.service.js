(function () {
    'use strict';
    angular.module('app')
        .factory('logBackend', logBackend);
    /* @ngInject */
    function logBackend($resource) {
        return {
            listPath: listPath,
            searchLogs: searchLogs,
            logContext: logContext
        };

        function listPath(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/search/paths/:app', {
                cluster: paramObj.cluster,
                user: paramObj.user,
                app: paramObj.app,
                task: paramObj.task,
                from: paramObj.from,
                to: paramObj.to
            });
        }

        function searchLogs(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/search/index', {
                cluster: paramObj.cluster,
                user: paramObj.user,
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
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/search/context', {
                cluster: paramObj.cluster,
                user: paramObj.user,
                app: paramObj.app,
                task: paramObj.task,
                path: paramObj.path,
                offset: paramObj.offset,
                page: paramObj.page,
                size: paramObj.size
            });
        }
    }
})();
