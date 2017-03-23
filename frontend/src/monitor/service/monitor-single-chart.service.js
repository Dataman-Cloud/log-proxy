(function () {
    'use strict';
    angular.module('app')
        .factory('monitorSingleChart', monitorSingleChart);


    /* @ngInject */
    function monitorSingleChart(chartUtil) {

        var optionsCls = createOptionsCls();
        var POINT_NUM = 180;

        return {
            Options: createOptions
        };

        function createOptions(type) {
            return new optionsCls(type);
        }

        function createOptionsCls(type) {
            function Options(type) {
                this.type = type;
                this.options = this._createOptions(this.type);
                this.data = [];
            }

            Options.prototype._createOptions = function (type) {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = type;
                options.chart.yAxis.tickFormat = function (d) {
                    return d;
                };
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype.pushData = function (data, api) {
                this._pushData(data.data.result, api, this.data);
            };


            Options.prototype._pushData = function (data, api, target) {
                angular.forEach(data, function (item, index) {
                    var serialKey = item.metric.container_env_mesos_task_id;
                    angular.forEach(item.values, function (value) {
                        chartUtil.pushData(target, serialKey, {
                            x: value[0] * 1000,
                            y: parseFloat(value[1])
                        }, POINT_NUM);
                    })

                }.bind(this))
            };

            Options.prototype.flushCharts = function (api) {
                api.update();
            };


            return Options;
        }
    }
})();
