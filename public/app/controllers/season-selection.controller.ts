'use strict';
import * as angular from 'angular';

/* Controllers */

angular.module('myApp.controllers', []).
  controller('SeasonSelectionCtrl', function ($scope, response) {
    $scope.data = response.data;
  });
