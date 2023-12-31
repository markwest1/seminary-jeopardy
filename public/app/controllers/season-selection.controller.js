'use strict';

/* Controllers */

angular.module('myApp.controllers', []).
  controller('SeasonSelectionCtrl', function ($scope, response) {
    $scope.data = response.data;
  });
