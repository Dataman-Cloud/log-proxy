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

            //about dashboard
            .state('home.dashboard', {
                url: '/dashboard',
                templateUrl: '/src/dashboard/dashboard.html',
                controller: 'DashboardCtrl as vm'
            })
            .state('home.dashboardMonitor', {
                url: '/dashboardMonitor/:clusterId',
                templateUrl: '/src/dashboard/monitor/app/list.html',
                controller: 'DashboardListAppCtrl as vm',
                resolve: {
                    info: getInfo
                }
            })
            .state('home.dashboardMonitor.detail', {
                url: '/:appId',
                templateUrl: '/src/dashboard/monitor/app/detail.html',
                controller: 'DashboardAppDetailCtrl as vm'
            })
            .state('home.dashboardInstanceMonitor', {
                url: '/dashboardMonitor/:clusterId/:appId/instances/:taskId',
                templateUrl: '/src/dashboard/monitor/instance/detail.html',
                controller: 'DashboardInstanceCtrl as vm'
            })
            //end dashboard

            //about monitor
            .state('home.monitor', {
                url: '/monitor?metric&appid&taskid&start&end&step&expr',
                templateUrl: '/src/monitor/monitorbase.html',
                controller: 'MonitorBaseCtrl as vm'
            })
            .state('home.monitor.chart', {
                url: '/chart',
                templateUrl: '/src/monitor/monitor/monitor.html',
                controller: 'MonitorCtrl as vm'
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
            .state('home.alertSilencesCreate', {
                url: '/alertSilencesCreate?fromByHistory',
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
        function getInfo(dashboardBackend, $stateParams) {
            return dashboardBackend.info({clusterid: $stateParams.clusterId}).get().$promise
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
            $state.go('home.dashboard');
        });
        $interpolateProvider.startSymbol('{/');
        $interpolateProvider.endSymbol('/}');
    }
})();
