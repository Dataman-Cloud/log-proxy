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
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/search/paths/:appid', {
                clusterid: paramObj.clusterid,
                userid: paramObj.userid,
                appid: paramObj.appid,
                taskid: paramObj.taskid,
                from: paramObj.from,
                to: paramObj.to
            });
        }

        function searchLogs(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/search/index', {
                clusterid: paramObj.clusterid,
                userid: paramObj.userid,
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
                clusterid: paramObj.clusterid,
                userid: paramObj.userid,
                appid: paramObj.appid,
                taskid: paramObj.taskid,
                path: paramObj.path,
                offset: paramObj.offset,
                page: paramObj.page,
                size: paramObj.size
            });
        }
    }
})();
