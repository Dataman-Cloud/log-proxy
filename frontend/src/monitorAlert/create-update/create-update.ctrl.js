/**
 * Created by lixiaoyan on 2017/3/24.
 */
(function () {
    'use strict';
    angular.module('app')
        .controller('CreateMonitorAlertCtrl', CreateMonitorAlertCtrl);
    /* @ngInject */

    function CreateMonitorAlertCtrl(target, monitorAlertBackend, Notification, $state, rule, alertBackend, $q) {
        var self = this;

        var ruleobj = rule.data || {};
        

        self.target = target;

        self.clusters = [];
        self.apps = [];
        self.indicators = [];
        self.unit = [];

        self.form = {
            class: ruleobj.class || '',
            name: ruleobj.name || '',
            cluster: ruleobj.cluster || '',
            app: ruleobj.app || '',
            severity: ruleobj.severity || '',
            indicator: ruleobj.indicator || '',
            pending: ruleobj.pending || '',
            aggregation: ruleobj.aggregation || '',
            comparison: ruleobj.comparison || '',
            threshold: ruleobj.threshold || '',
            cmdbAppid: '',
            description: ruleobj.description || ''

        };

        self.loadClusters = loadClusters;
        self.loadApps = loadApps;
        self.loadIndicators = loadIndicators;
        self.create = create;
        self.update = update;
        self.getCmdb = getCmdb;

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
            var obj = {
                appid: self.form.appid,
                cmdbAppid: self.form.cmdbAppid
            };
            $q.all({
                cmdbs: alertBackend.cmdbs().save(obj),
                alerts:monitorAlertBackend.rules().save(self.form)
            }).then(function (arr) {
                Notification.success('创建成功');
                $state.go('home.monitorAlert', null, {reload: true});
            });
        }

        // 更新
        function update() {
            self.form.ID = ruleobj.ID;
            monitorAlertBackend.rule(self.form.ID).update(self.form, function (data) {
                Notification.success('更新成功');
                $state.go('home.monitorAlert', null, {reload: true})
            });
        }

        function getCmdb() {
            alertBackend.cmdb(self.form.app).get(function (data) {
                self.form.cmdbAppid = data.data.cmdbAppid;
                self.form.appid = data.data.appid;
            })
        }
    }
})();
