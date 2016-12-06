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
            history: history,
            silences: silences,
            silence: silence
        };

        function alerts() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/alert', null, {
                'update': {method: 'PUT'}
            });
        }

        function alert(id) {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/alert/:id', {id: id});
        }

        function history() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/prometheus');
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