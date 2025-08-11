type FieldRelation = {
	type: 'relation';
	relation: {
		to: string;
	};
};

type FieldText = {
	type: 'text';
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

type Field = FieldRelation | FieldText | FieldDateTime | FieldDate | FieldNumber | FieldEmail;

export type CollectionDefinition = {
	slug: string;
	name: string;
	created_at: string;
	fields: Field[];
};
