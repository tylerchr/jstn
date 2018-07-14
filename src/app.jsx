import React, { Component, Fragment } from 'react';
import { Tokenizer, Parser, isValid } from 'jstn';

import ValidatedContent from './content.jsx';

import style from './app.less';

class App extends Component {

	constructor() {
		super();
		this.state = {
			jstnContent: "string?",
			jsonContent: "\"foo bar\"",
			parsedJSON: "foo bar",
			parsedJSTN: (new Parser("string?")).parseType(),
			validates: true
		};
	}

	componentDidMount() {
		console.log(Tokenizer, Parser, isValid); // eslint-disable-line
	}

	onChange(box, newValue) {
		this.setState({ [box]: newValue });

		if (box == "jsonContent") {
			try {
				let parsedJSON = JSON.parse(newValue);
				this.setState({ parsedJSON });
			} catch (e) {
				this.setState({ parsedJSON: undefined });
			}
		} else if (box == "jstnContent") {
			try {
				let parsedJSTN = (new Parser(newValue)).parseType();
				this.setState({ parsedJSTN });
			} catch (e) {
				this.setState({ parsedJSTN: undefined });
			}
		}

		var validates = false;
		if (this.state.parsedJSON !== undefined && this.state.parsedJSTN !== undefined) {
			validates = isValid(this.state.parsedJSTN, this.state.parsedJSON)
		}

		this.setState({ validates: validates })
	}

	onBlur(box) {
		if (box == "jsonContent" && this.state.parsedJSON != undefined) {
			this.setState({
				jsonContent: JSON.stringify(this.state.parsedJSON, null, "\t")
			});
		}
	}

	render() {

		// 'invalid-document'
		// 'valid-document'
		// 'everything-validates'

		let jstnValidation = (this.state.validates ? 'everything-validates' : (this.state.parsedJSTN === undefined ? 'invalid-document' : 'valid-document'));
		let jsonValidation = (this.state.validates ? 'everything-validates' : (this.state.parsedJSON === undefined ? 'invalid-document' : 'valid-document'));

		return (
			<Fragment>
				<header>
					<h1><strong>JSTN</strong><span>JSON Type Notation</span></h1>
				</header>
				<div className={style.wrapper}>
					<div className={style.banner}>
						<p>JSTN is a type declaration format for JSON documents. Writing a JSTN type declaration can help in:</p>
						<ul>
							<li>Describing APIs</li>
							<li>Validating JSON documents</li>
							<li>Enforcing type expectations</li>
						</ul>
						<p>Experiment below by writing a JSTN document and a JSON document.</p>
					</div>
					<ValidatedContent
						className={style.jstnEditor}
						title={'JSTN'}
						aceMode={'jstn'}
						value={this.state.jstnContent}
						valid={jstnValidation}
						onChange={this.onChange.bind(this, 'jstnContent')}
						onBlur={this.onBlur.bind(this, 'jstnContent')}
					/>
					<ValidatedContent
						className={style.jsonEditor}
						title={'JSON'}
						aceMode={'json'}
						value={this.state.jsonContent}
						valid={jsonValidation}
						onChange={this.onChange.bind(this, 'jsonContent')}
						onBlur={this.onBlur.bind(this, 'jsonContent')}
					/>
					<div className={style.main}>
						<p>TODO: More content here</p>
					</div>
				</div>
			</Fragment>
		)
	}
}

export default App;