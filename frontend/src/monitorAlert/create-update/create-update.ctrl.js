/**
 * Created by lixiaoyan on 2017/3/24.
 */
(function () {
    'use strict';
    angular.module('app')
        .controller('CreateMonitorAlertCtrl', CreateMonitorAlertCtrl);
    /* @ngInject */

    function CreateMonitorAlertCtrl(target, monitorAlertBackend, Notification, $state) {
        var self = this;
        self.target = target;
        self.clusters = [];
        self.apps = [];
        self.indicators = [];
        self.unit = [];
        self.form = {
            class: '',
            name: '',
            cluster: '',
            app: '',
            severity: '',
            indicator: '',
            pending: '',
            aggregation: '',
            comparison: '',
            threshold: ''
        };

        self.loadClusters = loadClusters;
        self.loadApps = loadApps;
        self.loadIndicators = loadIndicators;
        self.create = create;

        activate();

        function activate() {
            loadClusters();
            loadApps();
            loadIndicators()
        }

        // 集群
        function loadClusters() {
            monitorAlertBackend.clusters().get(function (data) {
                self.clusters = data.data;
            })
        }

        // 应用
        function loadApps() {
            monitorAlertBackend.apps().get(function (data) {
                self.apps = data.data;
            })
        }

        //指标
        function loadIndicators() {
            monitorAlertBackend.indicators().get(function (data) {
                angular.forEach(data.data, function (value, key) {
                    this.push(key);
                }, self.indicators);
                angular.forEach(data.data, function (value, key) {
                    this.push(value);
                }, self.unit);
            })
        }

        //创建
        function create() {
            monitorAlertBackend.rules().save(self.form, function (data) {
                Notification.success('创建成功');
                $state.go('home.monitorAlert', null, {reload: true})
            })
        }
    }
})();