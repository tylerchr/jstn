import React, { Component } from 'react';
import AceEditor from 'react-ace';

import 'brace';
import 'brace/mode/plain_text';
import 'brace/mode/json';
import 'brace/theme/github';

import style from './app.less';

class ValidatedContent extends Component {

	getValidationString(title, validationClass) {
		switch (validationClass) {
		case 'invalid-document':
			return 'Not a valid ' + title + ' document'
		case 'valid-document':
			return 'Valid ' + title + ' document'
		case 'everything-validates':
			if (title === 'JSTN') {
				return 'Describes the shape of the JSON document!';
			} else if (title === 'JSON') {
				return 'Valid with respect to the JSTN document!';
			}
		}

		return validationClass;
	}

	render() {

		let validationClass = {
			'everything-validates': 'validation-perfect',
			'valid-document': 'validation-valid',
			'invalid-document': 'validation-invalid',
		}[this.props.valid];

		return (
			<div className={style.editor + ' ' + this.props.className + ' ' + style[validationClass]}>
				<h2>{this.props.title}</h2>
				<div className={style.content}>
					<AceEditor
						mode={this.props.aceMode}
						height={'200px'}
						width={'474px'}
						showGutter={false}
						highlightActiveLine={false}
						wrapEnabled={true}
						fontSize={14}
						theme="github"
						value={this.props.value}
						onChange={this.props.onChange}
						onBlur={this.props.onBlur}
					/>
				</div>
				<div className={style.status}>
					{this.getValidationString(this.props.title, this.props.valid)}
				</div>
			</div>
		)
	}

}

export default ValidatedContent;