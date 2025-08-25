type BaseField = {
	name: string;
	label?: string;
	description?: string;
	// We don't have a choice yet for now.
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	constraints?: any;
};

type FieldSelect = {
	type: 'select';
};

type FieldRelation = {
	type: 'relation';
	relatesTo: string;
};

type FieldCheckbox = {
	type: 'checkbox';
};

type FieldRichText = {
	type: 'rich_text';
};

type FieldPlainText = {
	type: 'string';
};

type FieldDateTime = {
	type: 'date_time';
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

export type Collection = {
	name: string;
	slug: string;
};

export type Field = BaseField &
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
	updated_at: string;
	fields: {
		[key: string]: Field;
	};
};

export type CollectionResource = {
	id: string;
	slug: string;
	created_at: string;
	created_by: number;
	updated_at: string;
	updated_by: number;
	[key: string]: string | number | boolean | Date | undefined | null;
};
