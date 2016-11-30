(function () {
    'use strict';
    angular.module('app')
        .factory('logChart', logChart);


    /* @ngInject */
    function logChart($filter, chartUtil) {

        var optionsCls = createOptionsCls();
        var POINT_NUM = 12;

        return {
            Options: createOptions
        };

        function createOptions() {
            return new optionsCls();
        }

        function createOptionsCls() {
            function Options() {
                this.logOptions = this._createLogOptions();
                this.logData = [];
                this.dataChanged = false;
            }

            Options.prototype._createLogOptions = function () {
                var options = chartUtil.createDefaultOptions('discreteBarChart');
                options.title.text = '';
                options.chart.yAxis.tickFormat = function (d) {
                    return d;
                };
                options.chart.xAxis.tickFormat = function (d) {
                    return $filter('date')(d, 'M/d HH:mm');
                };
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype.pushData = function (data) {
                this._pushLogData(data);
                this.dataChanged = true;
            };

            Options.prototype._pushLogData = function (data) {
                this._pushData(data, this.logData);
            };

            Options.prototype._pushData = function (data, target) {
                angular.forEach(data, function (item) {
                    var serialKey;
                    serialKey = 'DocCount';
                    chartUtil.pushData(target, serialKey, {
                        x: item.key,
                        y: item.DocCount
                    }, POINT_NUM);

                })
            };

            return Options;
        }
    }
})();
