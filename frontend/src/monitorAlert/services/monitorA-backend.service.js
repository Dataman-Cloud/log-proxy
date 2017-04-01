/**
 * Created by lixiaoyan on 2017/3/24.
 */
(function () {
    'use strict';
    angular.module('app')
        .factory('monitorAlertBackend', monitorAlertBackend);

    /* @ngInject */

    function monitorAlertBackend($resource) {

        return {
            clusters: clusters,
            apps: apps,
            indicators: indicators,
            rules: rules,
            rule: rule,
            events: events,
            eventsAck: eventsAck
        };


        function clusters() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/clusters');
        }

        function apps(cluster) {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/clusters/:cluster/apps', {cluster: cluster});
        }

        function indicators() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/alert/indicators');
        }

        function rules() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/alert/rules');
        }

        function rule(id) {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/alert/rules/:id', {id: id}, {
                'update': {method: 'PUT'}
            });
        }

        function events(data) {
            var paramObj = data || {};

            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/alert/events', {
                cluster: paramObj.cluster,
                app: paramObj.app,
                ack: paramObj.ack,
                start: paramObj.start,
                end: paramObj.end,
                page: paramObj.page,
                size: paramObj.size
            });
        }

        function eventsAck(data) {
            var paramObj = data || {};

            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/alert/events/:id', {
                id: paramObj.id,
                cluster: paramObj.cluster,
                app: paramObj.app,
                ack: paramObj.ack,
                start: paramObj.start,
                end: paramObj.end,
                page: paramObj.page,
                size: paramObj.size,
                action: paramObj.action
            }, {
                'update': {method: 'PUT'}
            });
        }


    }
})();