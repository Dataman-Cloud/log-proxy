/**
 * Created by lixiaoyan on 2016/12/2.
 */
(function () {
    'use strict';
    angular.module('app')
        .factory('alertBackend', alertBackend);

    /* @ngInject */
    function alertBackend($resource) {
        return {
            alerts: alerts,
            alert: alert,
            histories: histories,
            clusters: clusters,
            apps: apps,
            silences: silences,
            silence: silence
        };

        function alerts(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/rules', {
                page: paramObj.page,
                size: paramObj.size
            }, {
                'update': {method: 'PUT'}
            });
        }

        function alert(id) {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/rules/:id', {id: id});
        }

        function histories(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/alerts', {
                cluster: paramObj.cluster,
                app: paramObj.app,
                task: paramObj.task,
                source: paramObj.source,
                start: paramObj.start,
                end: paramObj.end,
                page: paramObj.page,
                size: paramObj.size
            });
        }

        function clusters(data) {
            var paramObj = data || {};

            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/alerts/clusters', {
                start: paramObj.start,
                end: paramObj.end
            });
        }

        function apps(data) {
            var paramObj = data || {};

            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/log/alerts/clusters/:cluster/apps', {
                cluster: paramObj.cluster,
                start: paramObj.start,
                end: paramObj.end
            });
        }

        function silences() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/silences');
        }

        function silence(id) {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/silence/:id', {id: id}, {
                'update': {method: 'PUT'}
            });
        }
    }
})();
