CREATE SCHEMA backend_test;

create table if not exists backend_test.campaigns
(
    id    SERIAL primary key,
    title varchar(255)
);

create table if not exists backend_test.payment_methods
(
    id    SERIAL primary key,
    name varchar(255)
);

create table if not exists backend_test.donations
(
    id    SERIAL not null primary key,
    payment_method_id int,
    campaign_id int,
    amount decimal(9,2)
);

insert into backend_test.campaigns (id, title) values (0, 'campaign 1');
-- insert into backend_test.campaigns (id, title) values (2, 'campaign 2');

insert into backend_test.payment_methods (id, name) values (0, 'payment method 1');
-- insert into backend_test.payment_methods (id, name) values (2, 'payment method 2');

insert into backend_test.donations(id, payment_method_id, campaign_id, amount) values (0, 0, 0, 10000.00);