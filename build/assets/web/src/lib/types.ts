// Types mirroring the resume.json emitted by cmd/resume (internal/webdata).

export interface Socials {
	github?: string;
	gitlab?: string;
	linkedin?: string;
}

export interface Contact {
	firstname: string;
	lastname: string;
	title?: string;
	email?: string;
	phone?: string;
	location?: string;
	photo?: string;
	socials: Socials;
}

export interface ProjectLink {
	name: string;
	url: string;
	description?: string;
	icon?: string;
}

export interface PersonalProjects {
	intro?: string;
	links?: ProjectLink[];
}

export interface Entry {
	role: string;
	org?: string;
	orgUrl?: string;
	url?: string;
	location?: string;
	mode?: string;
	start?: string;
	end?: string;
	note?: string;
	intro?: string;
	bullets?: string[];
}

export interface Cert {
	title: string;
	org?: string;
	orgUrl?: string;
	certUrl?: string;
	location?: string;
	start?: string;
	end?: string;
	group?: string;
}

export interface SkillGroup {
	category: string;
	items: string[];
}

export interface Language {
	name: string;
	level?: string;
}

export interface Download {
	label: string;
	href: string;
	kind: 'PDF' | 'MD';
	group: 'Summary' | 'Detailed';
}

export interface Variant {
	experience: Entry[];
	education: Entry[];
	certificates: Cert[];
	activities: Entry[];
}

export interface Resume {
	contact: Contact;
	summary?: string;
	personalProjects: PersonalProjects;
	releasesUrl?: string;
	siteUrl?: string;
	downloads: Download[];
	skills: SkillGroup[];
	languages: Language[];
	variants: {
		summary: Variant;
		detailed: Variant;
	};
}

export type VariantName = 'summary' | 'detailed';
