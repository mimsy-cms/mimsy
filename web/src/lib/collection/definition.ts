type BaseField = {
	name: string;
	label: string;
};

type FieldSelect = {
	type: 'select';
};

type FieldRelation = {
	type: 'relation';
	relation: {
		to: string;
	};
};

type FieldCheckbox = {
	type: 'checkbox';
};

type FieldRichText = {
	type: 'richtext';
};

type FieldPlainText = {
	type: 'plaintext';
};

type FieldDateTime = {
	type: 'datetime';
};

type FieldDate = {
	type: 'date';
};

type FieldNumber = {
	type: 'number';
};

type FieldEmail = {
	type: 'email';
};

type Field = BaseField &
	(
		| FieldSelect
		| FieldRichText
		| FieldPlainText
		| FieldRelation
		| FieldCheckbox
		| FieldDateTime
		| FieldDate
		| FieldNumber
		| FieldEmail
	);

export type CollectionDefinition = {
	slug: string;
	name: string;
	created_at: string;
	fields: Field[];
};
