/**
 * Created by lixiaoyan on 2017/3/24.
 */
(function () {
    'use strict';
    angular.module('app')
        .controller('CreateMonitorAlertCtrl', CreateMonitorAlertCtrl);
    /* @ngInject */

    function CreateMonitorAlertCtrl(target, monitorAlertBackend, Notification, $state, rule) {
        var self = this;

        var rule = rule.data || {};

        self.target = target;

        self.clusters = [];
        self.apps = [];
        self.indicators = [];
        self.unit = [];
        
        self.form = {
            class: rule.class || '',
            name: rule.name || '',
            cluster: rule.cluster || '',
            app: rule.app || '',
            severity: rule.severity || '',
            indicator: rule.indicator || '',
            pending: rule.pending || '',
            aggregation: rule.aggregation || '',
            comparison: rule.comparison || '',
            threshold: rule.threshold || ''
        };

        self.loadClusters = loadClusters;
        self.loadApps = loadApps;
        self.loadIndicators = loadIndicators;
        self.create = create;
        self.update = update;

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
            monitorAlertBackend.apps(self.form.cluster).get(function (data) {
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

        // 更新
        function update() {
            self.form.ID = rule.ID;
            monitorAlertBackend.rule(self.form.ID).update(self.form, function (data) {
                Notification.success('更新成功');
                $state.go('home.monitorAlert', null, {reload: true})
            });
        }
    }
})();