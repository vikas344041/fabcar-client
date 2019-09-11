// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();
	
	$scope.queryAllCars = function(){

		appFactory.queryAllCars(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				data[i].Record.Id = data[i].Id;
				array.push(data[i].Record);
				console.log(array);
			}
			console.log(array);
			array.sort(function(a, b) {
			    return parseFloat(a.Id) - parseFloat(b.Id);
			});
			$scope.all_cars = array;
		});
	}

	$scope.queryCar = function(){

		var id = $scope.car_id;

		appFactory.queryCar(id, function(data){
			$scope.query_car = data;

			if ($scope.query_car == "Could not locate car"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	$scope.recordCar = function(){

		appFactory.recordCar($scope.car, function(data){
			$scope.create_car = data;
			$("#success_create").show();
		});
	}

	$scope.changeHolder = function(){

		appFactory.changeHolder($scope.holder, function(data){
			$scope.change_holder = data;
			if ($scope.change_holder == "Error: no car found"){
				$("#error_holder").show();
				$("#success_holder").hide();
			} else{
				$("#success_holder").show();
				$("#error_holder").hide();
			}
		});
	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllCars = function(callback){

    	$http.get('/get_all_cars/').success(function(output){
			callback(output)
		});
	}

	factory.queryCar = function(id, callback){
    	$http.get('/get_car/'+id).success(function(output){
			callback(output)
		});
	}

	factory.recordCar = function(data, callback){

		var car = data.id + "-" + data.model + "-" + data.make + "-" + data.owner + "-" + data.color;

    	$http.get('/add_car/'+car).success(function(output){
			callback(output)
		});
	}

	factory.changeHolder = function(data, callback){

		var holder = data.id + "-" + data.name;

    	$http.get('/change_holder/'+holder).success(function(output){
			callback(output)
		});
	}

	return factory;
});


