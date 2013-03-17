angular.module('vger', ['ui']).controller('tasks_ctrl',
	function ($scope, $http) {
		function get_process() {
			$http.get('/progress').success(function(data) {
				for (var i = data.length - 1; i >= 0; i--) {
					var item = data[i];
					item.StartDate = new Date(Date.parse(item.StartDate))
				};

				var collection  =  {};
				for (var i = data.length - 1; i >= 0; i--) {
					var item = data[i]
					collection[item.Name] = item;
				};
				var tasks = $scope.tasks;
				for (var i = tasks.length - 1; i >= 0; i--) {
					var task = tasks[i];
					var source = collection[task.Name];
					if (source) {
						angular.forEach(source, function (val, key) {
							task[key] = val;
						});
						delete collection[task.Name];
					} else {
						tasks.splice(i, 1);
					}
				}
				angular.forEach(collection, function (val, key) {
					tasks.push(val)
				});
			})
		}
		function init() {
			$http.get('/progress').success(function(data) {
				for (var i = data.length - 1; i >= 0; i--) {
					var item = data[i];
					item.StartDate = new Date(Date.parse(item.StartDate))
				};
				$scope.tasks = data;
			})
		}

		init();
		setInterval(get_process, 2000)

		$scope.parse_duration = function(dur) {
			var sec = Math.floor(dur/1000000000);
			var min = Math.floor(sec/60);
			var hour = Math.floor(min/60);
			var day = Math.floor(hour/24);

			return (day>0?day+'d':'')+(hour>0?hour%24+'h':'')+(min>0?min%60+'m':'')+(sec>0?sec%60+'s':'');
		}

		$scope.send_open = function (task) {
			$http.get('/open/' + task.Name).success(function (resp) {
				resp && $scope.push_alert(resp);
			});
		}
		$scope.send_resume = function (task) {
			$http.get('/resume/' + task.Name).success(function (resp) {
				resp && $scope.push_alert(resp);
				get_process();
			});
		}
		$scope.send_stop = function (task) {
			$http.get('/stop/' + task.Name).success(function (resp) {
				resp && $scope.push_alert(resp);
				get_process();
			});
		}
		$scope.send_limit = function (task) {
			$http.post('/limit/' + task.Name, task.LimitSpeed).success(function (resp) {
				resp && $scope.push_alert(resp);
				get_process();
			});
		};
		$scope.send_play = function (task) {
			$http.get('/play/' + task.Name).success(function (resp) {
				resp && $scope.push_alert(resp);
				get_process();
			})
		};

		$scope.waiting = false;
		function new_task() {
			$scope.waiting = true;
			if ($scope.new_url.indexOf('lixian.vip.xunlei.com') != -1) {
				$http.post('/new', $scope.new_url).success(function(resp) {
					$scope.new_url = '';
					$scope.waiting = false;
					resp && $scope.push_alert(resp);
				}).error(function(){$scope.waiting = false;});
			} else {
				$http.post('/thunder/new', $scope.new_url).success(function(data) {
					$scope.waiting = false;
					if (typeof data == 'string') {
						$scope.push_alert(data);
						return;
					}
					for (var i = data.length - 1; i >= 0; i--) {
						var item = data[i];
						item.loading = false;
					};
					$scope.bt_files = data;
				});
			}
		};
		$scope.get_bt_file_status = function(percent) {
			return (percent == 100) ? 'Finished' : percent + '%'
		}

		$scope.download_bt_files = function(file) {
			file.loading = true;

			$http.post('/new/' + file.Name, file.DownloadURL).success(
			function(resp) {
				file.loading = false;
				if (resp) $scope.push_alert(resp);
				else {
					$scope.bt_files = [];
					$scope.new_url = '';
					$scope.waiting = false;
				}
			}).error(function() {
				file.loading = false;
				$scope.waiting = false;
			});
		};
		$scope.bt_files = [];

		$scope.move_to_trash = function (task) {
			$http.get('/trash/' + task.Name).success(
			function (resp) {
				resp && $scope.push_alert(resp)
				get_process();
			});
		};

		$scope.set_autoshutdown = function (task) {
			$http.post('/autoshutdown/' + task.Name, task.Autoshutdown?'on':'off')
				.success(function(){
					get_process();
				});
		};


		//subtitles
		$scope.subtitles = [];
		$scope.subtitles_movie_name = '';

		$scope.search_subtitles = function(name) {
			$scope.subtitles_movie_name = name;
			$http.get('/subtitles/search/' + name).success(function (data) {
				if (data.length == 0) {
					$scope.nosubtitles = true;
					$scope.waiting = false;
					return
				}
				$scope.nosubtitles = false;
				for (var i = data.length - 1; i >= 0; i--) {
					var item = data[i];
					item.loading = false;

					//truncate description
					var description = item.Description;
					if (description.length > 50)
						item.Description = description.substr(0, 50) + '...';

				};
				$scope.subtitles = data;
				$scope.waiting = false;
			});
		};

		$scope.download_subtitles = function (sub) {
			sub.loading = true;
			$http.post('/subtitles/download/'+$scope.subtitles_movie_name, sub.URL).success(function () {
				sub.loading = false;
				$scope.subtitles = [];
			})
		};

		$scope.go = function () {
			$scope.waiting = true
			if (/.+\:\/\/.+|^magnet\:\?.+/.test($scope.new_url)) {
				new_task();
			} else {
				$scope.search_subtitles($scope.new_url)
			}
		};
		$scope.google_subtitles = function() {
			var name = $scope.subtitles_movie_name;
			name = name.replace(/(.*)[.](mkv|avi|mp4|rm|rmvb)/, '$1').replace(/(.*)-.*/, '$1') + ' subtitles';
			window.open("http://www.google.com/search?q=" + encodeURIComponent(name));
			$scope.nosubtitles = false;
		};

		$scope.parse_time = function (time) {
			var d = new Date(time * 1000);
			return d.format('ddd mmm dd')
		}
		$scope.upload_torrent = function ($event) {
			$event.preventDefault();
			
			if ($event.dataTransfer.files.length == 0) {
				var items = $event.dataTransfer.items;
				if (items.length > 0 && items[0].type=='text/plain') {
					items[0].getAsString(function(str){
						$scope.$apply(function() {
							$scope.new_url = str;
						})
					});
				}
				return;
			}

			if (/[.]torrent$/.test($event.dataTransfer.files[0].name)) {
				return;
			}

			var xhr = new XMLHttpRequest;
			var fd = new FormData();
			fd.append('torrent', $event.dataTransfer.files[0], 'torrent');

			$scope.waiting = true;

			xhr.open('POST', '/thunder/torrent');
			xhr.send(fd);
			xhr.onreadystatechange = function() {
				 if(this.readyState == this.DONE) {
				 	// if (!$scope.waiting)) return;
					$scope.$apply(function () {
					 	$scope.waiting = false;
					});

				    if(this.status == 200 && this.responseText != null) {
				    	var responseText = this.responseText;
				    	$scope.$apply(function () {
				    		if (responseText[0] != '[') {
				    			responseText && $scope.push_alert(responseText);
					    	} else {
						    	$scope.bt_files = JSON.parse(responseText);
						    }
				        });
				    }
				}
			}
		}

		$scope.alerts = [];
		$scope.push_alert = function(content, title) {
			title = title || 'Error';
			$scope.alerts.push({'title': title, 'content': content});
		}
		$scope.pop_alert = function() {
			$scope.alerts.pop();
		}
	}
);