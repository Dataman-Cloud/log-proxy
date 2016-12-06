(function () {
    'use strict';
    angular.module('app')
        .config(configure);

    /* @ngInject */
    function configure($locationProvider, $urlRouterProvider, $interpolateProvider, $stateProvider) {

        $stateProvider
            .state('home', {
                url: '/ui',
                templateUrl: '/src/app.html',
                controller: 'RootCtrl as vm',
                abstract: true
            })

            //about monitor
            .state('home.appmonitor', {
                url: '/appmonitor',
                templateUrl: '/src/monitor/app/list.html',
                controller: 'MonitorListAppCtrl as vm',
                resolve: {
                    apps: listApp
                }
            })
            .state('home.appmonitor.detail', {
                url: '/:appId',
                controller: 'MonitorAppDetailCtrl as vm',
                templateUrl: '/src/monitor/app/detail.html'
            })
            .state('home.instancemonitor', {
                url: '/appmonitor/:appId/instances/:taskId',
                templateUrl: '/src/monitor/instance/detail.html',
                controller: 'MonitorInstanceCtrl as vm'
            })
            .state('home.nodemonitor', {
                url: '/nodemonitor',
                templateUrl: '/src/monitor/node/node.html',
                controller: 'MonitorNodeCtrl as vm',
                resolve: {
                    nodes: listNode
                }
            })
            //end monitor

            //about log
            .state('home.logbase', {
                url: '/logbase',
                templateUrl: '/src/log/logbase.html',
                controller: 'LogBaseCtrl as vm'
            })
            .state('home.logbase.logs', {
                url: '/logs?appid&taskid&path&from&to&keyword',
                templateUrl: '/src/log/logs/logs.html',
                controller: 'LogsCtrl as vm'
            })
            .state('home.logbase.logWithoutKey', {
                url: '/logWithoutKey?appid&taskid&path&from&to',
                templateUrl: '/src/log/logs-without-keyword/logs.html',
                controller: 'LogWithoutKeyCtrl as vm'
            })
            .state('home.logbase.logcontext', {
                url: '/logcontext?appid&taskid&path&offset',
                templateUrl: '/src/log/log-context/logs.html',
                controller: 'LogContextCtrl as vm'
            })
            //end log

            //about alert
            .state('home.alertKeywordCreate', {
                url: '/alertKeywordCreate',
                templateUrl: '/src/alert/keywords/create-update/create-update.html',
                controller: 'CreateAlertKeywordCtrl as vm',
                resolve: {
                    target: function () {
                        return 'create'
                    },
                    alert: function () {
                        return {}
                    }
                }
            })
            .state('home.alertKeywordUpdate', {
                url: '/alertKeywordUpdate/:id',
                templateUrl: '/src/alert/keywords/create-update/create-update.html',
                controller: 'CreateAlertKeywordCtrl as vm',
                resolve: {
                    target: function () {
                        return 'update'
                    },
                    alert: getAlert
                }
            })
            .state('home.alertkeyword', {
                url: '/alertkeyword',
                templateUrl: '/src/alert/keywords/list/list.html',
                controller: 'AlertKeywordsCtrl as vm'
            })
            .state('home.alerthistory', {
                url: '/alerthistory',
                templateUrl: '/src/alert/history/list/list.html',
                controller: 'AlertHistoriesCtrl as vm'
            })
            .state('home.alertSilences', {
                url: '/alertSilences',
                templateUrl: '/src/alert/silence/list/list.html',
                controller: 'AlertSilencesCtrl as vm'
            })
            .state('home.historySilence', {
                url: '/historySilence',
                templateUrl: '/src/alert/history/silence/silence.html',
                controller: 'AlertSilencesCtrl as vm'
            })
            .state('home.alertSilencesCreate', {
                url: '/alertSilencesCreate',
                templateUrl: '/src/alert/silence/create-update/create-update.html',
                controller: 'CreateAlertSilenceCtrl as vm',
                resolve: {
                    target: function () {
                        return 'create'
                    },
                    silence: function () {
                        return {}
                    }
                }
            })
            .state('home.alertSilencesUpdate', {
                url: '/alertSilencesUpdate/:id?from',
                templateUrl: '/src/alert/silence/create-update/create-update.html',
                controller: 'CreateAlertSilenceCtrl as vm',
                resolve: {
                    target: function () {
                        return 'update'
                    },
                    silence: getSilence
                }
            });

        /* @ngInject */
        function listApp(monitorBackend) {
            return monitorBackend.listApp().get().$promise
        }

        /* @ngInject */
        function listNode(monitorBackend) {
            return monitorBackend.listNode().get().$promise
        }

        /* @ngInject */
        function getAlert(alertBackend, $stateParams) {
            return alertBackend.alert($stateParams.id).get().$promise
        }

        function getSilence(alertBackend, $stateParams) {
            return alertBackend.silence($stateParams.id).get().$promise
        }

        $locationProvider.html5Mode(true);
        //warning: otherwise(url) will be redirect loop on state with errored resolve
        $urlRouterProvider.otherwise(function ($injector) {
            var $state = $injector.get('$state');
            $state.go('home.appmonitor');
        });
        $interpolateProvider.startSymbol('{/');
        $interpolateProvider.endSymbol('/}');
    }
})();
