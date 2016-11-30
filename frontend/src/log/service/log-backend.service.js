(function () {
    'use strict';
    angular.module('app')
        .factory('logBackend', logBackend);
    /* @ngInject */
    function logBackend($resource) {
        return {
            listApp: listApp,
            listTask: listTask,
            listPath: listPath,
            searchLogs: searchLogs,
            logContext: logContext
        };

        function listApp(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/search/applications', {
                from: paramObj.from,
                to: paramObj.to
            });
        }

        function listTask(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/search/tasks/:appid', {
                appid: paramObj.appid,
                from: paramObj.from,
                to: paramObj.to
            });
        }

        function listPath(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/search/paths/:appid/:taskid', {
                appid: paramObj.appid,
                taskid: paramObj.taskid,
                from: paramObj.from,
                to: paramObj.to
            });
        }

        function searchLogs(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/search/index', {
                appid: paramObj.appid,
                taskid: paramObj.taskid,
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
                appid: paramObj.appid,
                taskid: paramObj.taskid,
                path: paramObj.path,
                offset: paramObj.offset
            });
        }
    }
})();
