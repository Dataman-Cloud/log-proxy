(function () {
    'use strict';
    angular.module('app')
        .controller('CreateAlertKeywordCtrl', CreateAlertKeywordCtrl);
    /* @ngInject */
    function CreateAlertKeywordCtrl(alertBackend, $state, target, alert, Notification, logBackend, $q) {
        var alert = alert.data || {};

        var self = this;

        self.target = target;
        self.form = {
            cluster: alert.cluster,
            app: alert.app,
            source: alert.source || '',
            keyword: alert.keyword || '',
            cmdbAppid: '',
            description: alert.description || ''

        };

        self.clusters = {};
        self.apps = {};
        self.sources = {};

        self.loadClusters = loadClusters;
        self.loadApps = loadApps;
        self.loadSources = loadSources;
        self.create = create;
        self.update = update;
        self.getCmdb = getCmdb;

        activate();

        function activate() {
            loadClusters();
            loadApps();
            loadSources();
        }

        function loadClusters() {
            logBackend.clusters().get(function (data) {
                self.clusters = data.data;
            })
        }

        function loadApps() {
            if (self.form.cluster) {
                logBackend.apps(self.form.cluster).get(function (data) {
                    self.apps = data.data;
                });
            }
        }

        function loadSources() {
            if (self.form.app && self.form.cluster) {
                logBackend.sources({
                    cluster: self.form.cluster,
                    app: self.form.app
                }).get(function (data) {
                    self.sources = data.data;
                })
            }
        }

        function create() {
            var obj = {
                appid: self.form.appid,
                cmdbAppid: self.form.cmdbAppid
            };
            $q.all({
                cmdbs: alertBackend.cmdbs().save(obj),
                alerts:alertBackend.alerts().save(self.form)
            }).then(function (arr) {
                Notification.success('创建成功');
                $state.go('home.alertkeyword', null, {reload: true});
                
            });
        }

        function update() {
            self.form.id = alert.id;
            alertBackend.alerts().update(self.form, function (data) {
                Notification.success('更新成功');
                $state.go('home.alertkeyword', null, {reload: true})
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
