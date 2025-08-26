export type User = {
	id: number;
	email: string;
	created_at: string;
	updated_at: string;
};

export type AuthUser = {
	id: number;
	email: string;
	must_change_password: boolean;
};
