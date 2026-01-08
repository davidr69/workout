const months = [null, 'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

export default class Render {
	monthList;

	init(monthList) {
		this.monthList = monthList['dates'];
	}

	getMonths() {
		return this.monthList;
	}

	column(yrmon, data) {
		let col = '#col' + yrmon;
		let nl = document.querySelectorAll(col);
		let nlBefore;
		let nlAfter;
		if(nl.length === 0) {
			let pos = this.monthList.indexOf(yrmon);
			if(pos > 0) {
				nlBefore = document.querySelectorAll('#col' + this.monthList[pos - 1]);
			}
			if(pos < this.monthList.length) {
				nlAfter = document.querySelectorAll('#col' + this.monthList[pos + 1]);
			}
			// all information has been gathered; now need to make a decision
			if(/*nlBefore.length == 0 &&*/ nlAfter.length === 0) {
				// just a simple add
				this.#append_column(yrmon, data); // offset from muscle category is 0
			}
		} else {
			//
		}
	}

	#append_column(yrmon, data) {
		let tr = document.getElementById('tableRecord');
		let th = document.createElement('th');
		let cellId = 'col' + yrmon;
		th.setAttribute('id', cellId);
		let year = yrmon.substring(0, 4);
		let month = yrmon.substring(4);
		let textNode = document.createTextNode(months[Number(month)] + ' ' + year);
		th.appendChild(textNode);
		tr.appendChild(th);

		let idx = 0;
		for(let el of document.querySelectorAll('.exercise')) {
			// find all the left-pane exercises and draw to their right
			let row = data[idx++];
			// row should have {"id":93,"muscle":"Abdominals","muscleId":8,"exercise":"Ab Carver","weight":null,"rep1":null,"rep2":null,"progId":123}
			let td = document.createElement('td');
			td.setAttribute('id', 'col' + yrmon);
			if(row['weight'] === null && row['rep1'] === null && row['rep2'] === null) {
				td.innerHTML = '&nbsp;';
			} else {
				let anchor = document.createElement('a');
				anchor.setAttribute('href', 'javascript:globalThis.workout.edit(' + row['progress_id'] + ')');

				let disp;
				if(row['weight'] === null) {
					disp = row['rep1'];
					if(row['rep2'] !== null) {
						disp += ', ' + row['rep2'];
					}
				} else {
					if(row['rep2'] === null) {
						disp = `${row['weight']} / ${row['rep1']}`;
					} else {
						disp = `${row['rep1']}, ${row['rep2']} (${row['weight']})`;
					}
				}

				let text = document.createTextNode(disp);
				td.setAttribute('class', 'center');

				anchor.appendChild(text);
				td.appendChild(anchor);
			}

			el.parentNode.appendChild(td);
		}

	}

	drawMusclesAndExercises(data) {
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
		let th = document.getElementById('tableBody').parentNode;
		data['exercises'].forEach(section => {
			let tr = document.createElement('tr');
			let td = document.createElement('td');
			td.setAttribute('class', 'muscle');
			let textNode = document.createTextNode(section['description']);
			td.appendChild(textNode);
			tr.appendChild(td);
			th.appendChild(tr);

			section['exercises'].forEach(obj => {
				let tr = document.createElement('tr');
				tr.setAttribute('class', 'data');

				let td = document.createElement('td');
				td.setAttribute('class', 'exercise');
				td.setAttribute('id', 'ex' + obj['id']);

				let textNode = document.createTextNode(obj['exercise_name']);
				td.appendChild(textNode);
				tr.appendChild(td);

				th.appendChild(tr);
			});
		});
	}
}
