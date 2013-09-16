angular.module('vger', ['ngAnimate', 'ui']).controller('tasks_ctrl',
	function($scope, $http) {

		function monitor(path, ondata, onclose, onerror) {
			var websocket = new WebSocket('ws://' +
				window.location.host + path);

			websocket.onopen = onOpen;
			websocket.onclose = onClose;
			websocket.onmessage = onMessage;
			websocket.onerror = onError;

			function onOpen(evt) {}

			function onClose(evt) {
				// $scope.push_alert('socket close');
				if (onclose) {
					onclose(evt);
				}
			}

			function onMessage(evt) {
				$scope.$apply(function() {
					ondata(JSON.parse(evt.data));
				});
			}

			function onError(evt) {
				// $scope.push_alert(evt.data);
				if (onerror) {
					onerror(evt.data);
				}
			}

			function doSend(message) {
				websocket.send(message);
			}
		}


		$scope.tasks = [];
		$scope.config = {
			'max-speed': '0'
		};
		$http.get('/config').success(function(resp) {
			$scope.config = resp;
			var v = $scope.config['shutdown-after-finish'];
			$scope.config['shutdown-after-finish'] = (v == 'true');
		})

		function monitor_process() {
			monitor('/progress', function(data) {
				for (var i = data.length - 1; i >= 0; i--) {
					var item = data[i];
					item.StartDate = new Date(Date.parse(item.StartDate))
				};

				var collection = {};
				for (var i = data.length - 1; i >= 0; i--) {
					var item = data[i]
					if (item.Status != "Deleted") {
						collection[item.Name] = item;
					}
				};
				var tasks = $scope.tasks;
				for (var i = tasks.length - 1; i >= 0; i--) {
					var task = tasks[i];
					var source = collection[task.Name];
					if (source) {
						angular.forEach(source, function(val, key) {
							task[key] = val;
						});
						delete collection[task.Name];
					} else {
						tasks.splice(i, 1);
					}
				}
				angular.forEach(collection, function(val, key) {
					tasks.push(val)
				});

			}, monitor_process);
		}
		monitor_process();
		$scope.new_url = document.getElementById('new-url').value;


		$scope.parse_duration = function(dur) {
			var sec = Math.floor(dur / 1000000000);
			var min = Math.floor(sec / 60);
			var hour = Math.floor(min / 60);
			var day = Math.floor(hour / 24);

			return (day > 0 ? day + 'd' : '') + (hour > 0 ? hour % 24 + 'h' : '') + (min > 0 ? min % 60 + 'm' : '') + (sec > 0 ? sec % 60 + 's' : '');
		}

		$scope.send_open = function(task) {
			$http.get('/open/' + task.Name).success(function(resp) {
				resp && $scope.push_alert(resp);
			});
		}
		$scope.send_resume = function(task) {
			$http.get('/resume/' + task.Name).success(function(resp) {
				resp && $scope.push_alert(resp);
			});
		}
		$scope.send_stop = function(task) {
			$http.get('/stop/' + task.Name).success(function(resp) {
				resp && $scope.push_alert(resp);
			});
		}
		$scope.send_limit = function($event) {
			$http.get('/limit/' + $event.target.value).success(function(resp) {
				resp && $scope.push_alert(resp);
			});
		};
		$scope.send_simultaneous_downloads = function() {
			$http.post('/config/simultaneous', $scope.config['simultaneous-downloads'])
				.success(function () {});
		}
		$scope.send_play = function(task) {
			$http.get('/play/' + task.Name).success(function(resp) {
				resp && $scope.push_alert(resp);
			})
		};

		$scope.waiting = false;

		function new_task() {
			$scope.waiting = true;
			if ($scope.new_url.indexOf('lixian.vip.xunlei.com') != -1 ||
				$scope.new_url.indexOf('youtube.com') != -1 ||
				/.*dmg|.*zip|.*rar|.*exe|.*iso/.test($scope.new_url)) {
				$http.post('/new/', $scope.new_url).success(function(resp) {
					if (!resp) {
						$scope.new_url = '';
					}
					$scope.waiting = false;
					resp && $scope.push_alert(resp);
				}).error(function() {
					$scope.waiting = false;
				});
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

						var j = item.Name.lastIndexOf('\/');
						item.Name = item.Name.substring(j + 1);
					}
					if (data.length == 1 && data[0].Percent == 100) {
						$scope.waiting = true;
						$scope.download_bt_files(data[0]);
					} else {
						$scope.bt_files = data;
					}
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
					$scope.waiting = false;
					if (resp) $scope.push_alert(resp);
					else {
						// $scope.bt_files = [];
						$scope.new_url = '';
					}
				}).error(function() {
				file.loading = false;
				$scope.waiting = false;
			});
		};
		$scope.bt_files = [];

		$scope.move_to_trash = function(task) {
			$http.get('/trash/' + task.Name).success(
				function(resp) {
					resp && $scope.push_alert(resp)
				});
		};

		$scope.set_autoshutdown = function() {
			$http.post('/autoshutdown', $scope.config['shutdown-after-finish']?'true':'false')
				.success(function() {});
		};


		//subtitles
		$scope.subtitles = [];
		$scope.subtitles_movie_name = '';

		$scope.search_subtitles = function(name) {
			$scope.subtitles = [];
			$scope.subtitles_movie_name = name;

			monitor('/subtitles/search/' + name, function(data) {
				$scope.nosubtitles = false;
				data.loading = false;

				//truncate description
				data.FullDescription = data.Description;
				var description = data.Description;
				if (description.length > 73)
					data.Description = description.substr(0, 35) + '...' + description.substr(description.length - 35, 35);

				$scope.subtitles.push(data);
				$scope.waiting = true;
			}, function() {
				if ($scope.subtitles.length == 0) {
					$scope.nosubtitles = true;
				}
				$scope.waiting = false;
			}, function() {
				if ($scope.subtitles.length == 0) {
					$scope.nosubtitles = true;
				}
				$scope.waiting = false;
			});
		};

		$scope.download_subtitles = function(sub) {
			sub.loading = true;
			$http.post('/subtitles/download/' + $scope.subtitles_movie_name, sub.URL).success(function() {
				sub.loading = false;
				$scope.subtitles = [];
			})
		};

		$scope.go = function() {
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
		$scope.addic7ed_subtitles = function() {
			var name = $scope.subtitles_movie_name;
			var i = name.lastIndexOf('.');
			if (i != -1) {
				name = name.substring(0, i);
			}

			i = name.lastIndexOf('-');
			if (i != -1) {
				name = name.substring(0, i);
			}

			name = name.replace(/720p|x[.]264|BluRay|DTS|x264|1080p|H[.]264|AC3|[.]ENG|[.]BD|[.]Rip|H264|HDTV|-IMMERSE|-DIMENSION|xvid|[[]PublicHD[]]|[.]Rus|Chi_Eng|DD5[.]1|HR-HDTV|[.]HDTVrip|[.]AAC|[0-9]+x[0-9]+|blu-ray|Remux|dxva|dvdscr|WEB-DL/ig, '');
			name = name.replace(/[\u4E00-\u9FFF]+/ig, '');
			name = name.replace(/[.]/g, ' ');

			window.open("http://www.addic7ed.com/search.php?search=" + encodeURIComponent(name));
			$scope.nosubtitles = false;
		}


		$scope.parse_time = function(time) {
			var d = new Date(time * 1000);
			return d.format('ddd mmm dd')
		}
		$scope.upload_torrent = function($event) {
			$event.preventDefault();

			if ($event.dataTransfer.files.length == 0) {
				var items = $event.dataTransfer.items;
				if (items && items.length > 0 && items[0].type == 'text/plain') {
					items[0].getAsString(function(str) {
						$scope.$apply(function() {
							$scope.new_url = str;
						})
					});
				}
				return;
			}

			if (!/[.]torrent$/.test($event.dataTransfer.files[0].name)) {
				$scope.push_alert('Only support .torrent file!')
				return;
			}

			var xhr = new XMLHttpRequest;
			var fd = new FormData();
			fd.append('torrent', $event.dataTransfer.files[0], 'torrent');

			$scope.waiting = true;

			xhr.open('POST', '/thunder/torrent');
			xhr.send(fd);
			xhr.onreadystatechange = function() {
				if (this.readyState == this.DONE) {
					// if (!$scope.waiting)) return;
					$scope.$apply(function() {
						$scope.waiting = false;
					});

					if (this.status == 200 && this.responseText != null) {
						var responseText = this.responseText;
						$scope.$apply(function() {
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
			$scope.alerts.push({
				'title': title,
				'content': content
			});
		}
		$scope.pop_alert = function() {
			$scope.alerts.pop();
		}

		window.onload = function() {
			setTimeout(function() {
				document.getElementById('box-overlay').style.display = '';
			}, 500);
			var ele = document.getElementById('new-url');
			ele.value = getCookie('input');
			ele.select();

			$scope.new_url = ele.value;
		}
	}

);


function getCookie(c_name) {
	var c_value = document.cookie;
	var c_start = c_value.indexOf(" " + c_name + "=");
	if (c_start == -1) {
		c_start = c_value.indexOf(c_name + "=");
	}
	if (c_start == -1) {
		c_value = null;
	} else {
		c_start = c_value.indexOf("=", c_start) + 1;
		var c_end = c_value.indexOf(";", c_start);
		if (c_end == -1) {
			c_end = c_value.length;
		}
		c_value = unescape(c_value.substring(c_start, c_end));
	}
	return c_value;
}

function setCookie(c_name, value, exdays) {
	var exdate = new Date();
	exdate.setDate(exdate.getDate() + exdays);
	var c_value = escape(value) + ((exdays == null) ? "" : "; expires=" + exdate.toUTCString());
	document.cookie = c_name + "=" + c_value;
}
window.onbeforeunload = function() {
	setCookie('input', document.getElementById('new-url').value, 10000)
}