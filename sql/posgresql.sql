create table users (
	user_id serial primary key,
	name text not null,
	lastname text not null,
	sex int,
	sex_text text generated always as (case when sex = 0 then 'мужской' else 'женский' end) stored,
	age int not null,
	is_deleted bool not null default false
);

create table segments (
	segment_id serial primary key,
	name text not null,
	is_deleted bool not null default false,
	unique (name)
);

create table users_segments (
	user_segment_id serial primary key,
	user_id int not null,
	segment_id int not null,
	start_date timestamp not null default current_timestamp,
	end_date timestamp,
	foreign key (user_id) references users (user_id) on delete no action,
	foreign key (segment_id) references segments(segment_id) on delete no action
);

insert into users (name, lastname, sex, age)
values
    ('Иван', 'Иванов', '0', 32),
    ('Анна', 'Петрова', '1', 25),
    ('Алексей', 'Сидоров', '0', 42),
    ('Мария', 'Федорова', '1', 19),
    ('Петр', 'Николаев', '0', 57),
    ('Екатерина', 'Кузнецова', '1', 38),
    ('Андрей', 'Михайлов', '0', 41),
    ('Ольга', 'Андреева', '1', 29),
    ('Дмитрий', 'Козлов', '0', 24),
    ('Маргарита', 'Волкова', '1', 30),
    ('Сергей', 'Захаров', '0', 47),
    ('Елена', 'Сергеева', '1', 36),
    ('Никита', 'Романов', '0', 22),
    ('Виктория', 'Ильина', '1', 27),
    ('Александр', 'Гаврилов', '0', 39);
	
insert into segments (name)
values
    ('AVITO_VOICE_MESSAGES'),
    ('AVITO_MARKET'),
    ('AVITO_DELIVERY'),
    ('AVITO_DISCOUNT_30'),
    ('AVITO_DISCOUNT_50'),
    ('AVITO_DISCOUNT_70'),
    ('AVITO_MARKET_DISCOUNT_30'),
    ('AVITO_MARKET_DISCOUNT_45');
	
insert into users_segments (user_id, segment_id, start_date, end_date)
values
    (1, 1, '2021-01-01T12:35:50', null),
    (1, 2, '2021-01-01T12:35:50', '2022-01-01'),
    (1, 3, '2021-02-01T12:35:50', '2024-08-01'),
    (2, 1, '2021-03-01T12:35:46', null),
    (2, 2, '2021-05-01T12:35:46', '2022-01-01'),
    (4, 5, '2021-06-01T16:10:22', null),
    (4, 6, '2021-07-01T16:10:22', null),
    (5, 1, '2021-09-01T10:00:25', null);
