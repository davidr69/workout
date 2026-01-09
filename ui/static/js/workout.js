/**
 * The steps that need to be performed are:
 * 1. get a list of all exercises
 * 2. retrieve a list of all month/years
 *
 * After BOTH of those operations complete, the statistics can be populated:
 * 3. retrieve that status for the last n columns iteratively and asynchronously
 *
 * After n column stats are retrieved, rendering can begin. These cannot be done
 * in parallel since the insertion/append point is determined once, and parallel
 * rendering can cause columns to switch progress, changing the insert/append point.
 * When a column render is requested, the rendering class will determine if it needs
 * to be retrieved or if it simply needs to be made visible (previously drawn). If
 * the data needs to be retrieved, the rendering class will determine where it needs
 * to be placed in the DOM.
 * 4. Render n columns iteratively and synchronously
 */
import Render from "./render.js";

export default class Workout {
	render = new Render();
	allData = {};
	count = 0;
	monthCount;
	exercises;
	headers = {
		'Content-Type': 'application/json',
		'Accept': 'application/json'
	};

	constructor() {
		this.monthCount = navigator.userAgent.search('Mobile') === -1 ? 6 : 3;

		this.#drawTable();

		const formatter = new Intl.DateTimeFormat('en', { month: 'long' });

		const months = Array.from({ length: 12 }, (_, i) =>
			formatter.format(new Date(2000, i, 1)));

		const today = new Date();
		const currentMonth = today.getMonth();
		const monthDropDown = document.getElementById("choose_month");

		for(let month = 0; month < 12; month++) {
			let option = new Option(months[month], String(month + 1));
			monthDropDown.add(option);
		}
		monthDropDown.selectedIndex = currentMonth;
		document.getElementById("choose_year").value = today.getFullYear();
	}

	#getMonths = () => {
 		fetch('api/v1/months').then(response => {
			response.json().then(data => {
				this.render.init(data);
				this.#drawHeaders();
			});
		});
    }

	#getHeader = (when, cb) => {
		const year = when.substring(0,4);
		const month = when.substring(4);
		fetch(`api/v1/progress?year=${year}&month=${month}`).then(response => {
			response.json().then(data => {
				this.allData[when] = data['progress'];
				if(--this.count === 0) {
					cb();
				}
			});
		});
	}

	#drawHeaders = () => {
		let l = this.render.getMonths();
		this.count = this.monthCount;
		let callback = () => {
			return this.#drawColumns();
		};
		for(let i = l.length - this.monthCount; i<l.length; i++) {
			this.#getHeader(l[i], callback);
		}
	}

	#drawColumns = () => {
		let l = this.render.getMonths();
		for(let i = l.length - this.monthCount; i<l.length; i++) {
			let yrmon = l[i];
			this.render.column(yrmon, this.allData[yrmon]);
		}
	}

	#drawTable = () => {
		// get list of all exercises (not activity)
		fetch('api/v1/exercises').then(response => {
			/*
				Data looks like:

				{
					"exercises": [
						{
							"muscle": "Abdominals",
							"exercises": [
								{
									"id": 93,
									"muscle": "Abdominals",
									"exercise_name": "Ab Carver"
								}
							]
						}
					]
				}
			*/
			response.json().then(data => {
				// draw the left side muscle/category list
				this.render.drawMusclesAndExercises(data);
				this.#getMonths();
				this.exercises = data['exercises'];
				this.#populateMuscles()
			});
		});
	}

	#populateMuscles = () => {
		let dd = document.getElementById('choose_muscle');
		dd.length = 0;
		let muscles = this.exercises.map(item => item.muscle).sort();
		for(let muscle of muscles) {
			let option = new Option(muscle, muscle);
			dd.add(option);
		}
		this.populateExercises();
	}

	populateExercises = () => {
		let dd = document.getElementById('choose_exercise');
		dd.length = 0;
		let muscle = document.getElementById('choose_muscle').value;
		let exercises = this.exercises.find(item => item.muscle === muscle).exercises;
		let data = exercises.map(item => [item['id'], item['exercise_name']] );
		for (let datum of data) {
			let option = new Option(datum[1], datum[0]);
			dd.add(option);
		}
	}

	rewind = () => {
		//
	}

	forward = () => {
		//
	}

	edit = (id) => {
		let url = `edit.html?id=${id}`;
		window.open(url, "Edit/Delete", "width=500,height=350");
	}

	#getFormValues = () => {
		const exercise = document.getElementById('choose_exercise').value;
		const month = document.getElementById('choose_month').value;
		const year = document.getElementById('choose_year').value;
		const weight = document.getElementById('weight').value.trim();
		const rep1 = document.getElementById('rep1').value.trim();
		const rep2 = document.getElementById('rep2').value.trim();

		const mydate = new Date(year, month - 1, 1);

		return {
			exercise: Number(exercise),
			mydate: mydate,
			weight: weight === '' ? null : Number(weight),
			rep1: rep1 === '' ? null : Number(rep1),
			rep2: rep2 === '' ? null : Number(rep2)
		};
	}

	save = () => {
		const body = this.#getFormValues();

		fetch('api/v1/activity', {
			method: 'POST',
			body: JSON.stringify(body),
			headers: this.headers
		}).then(response => {
			response.json().then(data => {
				console.log(data);
				globalThis.location.reload();
			})
		}).catch(error => {
			console.log(error);
		});
	}
};
