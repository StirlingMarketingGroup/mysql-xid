# MySQL xid - Globally Unique ID Generator

`xid` is a cool library for generating 12 byte safe-from-anywhere IDs. This is a repo that wraps `xid` for use in MySQL as a MySQL UDF. The following is straight from the `xid` readme.

---

Package xid is a globally unique id generator library, ready to safely be used directly in your server code.

Xid uses the Mongo Object ID algorithm to generate globally unique ids with a different serialization (base64) to make it shorter when transported as a string:
https://docs.mongodb.org/manual/reference/object-id/

- 4-byte value representing the seconds since the Unix epoch,
- 3-byte machine identifier,
- 2-byte process id, and
- 3-byte counter, starting with a random value.

The binary representation of the id is compatible with Mongo 12 bytes Object IDs.
The string representation is using base32 hex (w/o padding) for better space efficiency
when stored in that form (20 bytes). The hex variant of base32 is used to retain the
sortable property of the id.

Xid doesn't use base64 because case sensitivity and the 2 non alphanum chars may be an
issue when transported as a string between various systems. Base36 wasn't retained either
because 1/ it's not standard 2/ the resulting size is not predictable (not bit aligned)
and 3/ it would not remain sortable. To validate a base32 `xid`, expect a 20 chars long,
all lowercase sequence of `a` to `v` letters and `0` to `9` numbers (`[0-9a-v]{20}`).

## Features:

- Size: 12 bytes (96 bits), smaller than UUID, larger than snowflake
- Base32 hex encoded by default (20 chars when transported as printable string, still sortable)
- Non configured, you don't need set a unique machine and/or data center id
- K-ordered
- Embedded time with 1 second precision
- Unicity guaranteed for 16,777,216 (24 bits) unique ids per second and per host/process
- Lock-free (i.e.: unlike UUIDv1 and v2)

## Notes:

- Xid is dependent on the system time, a monotonic counter and so is not cryptographically secure. If unpredictability of IDs is important, you should not use Xids. It is worth noting that most other UUID-like implementations are also not cryptographically secure. You should use libraries that rely on cryptographically secure sources (like /dev/urandom on unix, crypto/rand in golang), if you want a truly random ID generator.


---
## It's MySQL time

### `xid_bin`

Returns a new xid in the 12 byte binary version. Sadly this couldn't be called simply "xid", because I guess MySQL already has a native function with that name that is unrelated to this.

```sql
`xid` ( ) : binary(12)
```
---
### `xid_string`

Returns a new xid in the 20 char, base32 version (`[0-9a-v]{20}`).

```sql
`xid_string` ( ) : char(20)
```
---
### `xid_to_bin`

Takes an encoded xid and returns the binary version. Will return `null` if the value given was not a valid xid.

```sql
`xid_to_bin` ( char(20) `id` ) : binary(12)
```
### `bin_to_xid`

Takes an binary xid and returns the encoded version. Will return `null` if the value given was not a valid xid.

```sql
`bin_to_xid` ( binary(12) `id` ) : char(20)
```
---
## Examples

```sql
select`xid_bin`();
-- 0x604B7B580C43999B8B09D268

select`xid_string`();
-- 'c15nmq0c8ecpn2o9q9kg'

select`bin_to_xid`(0x604B7B580C43999B8B09D268);
-- 'c15nmm0c8ecpn2o9q9k0'
select`bin_to_xid`(null);
-- null
select`bin_to_xid`('yeet');
-- notice how xid doesn't throw errors for this input!
-- 'f5imat00000000000000'

select`xid_to_bin`('c15nmq0c8ecpn2o9q9kg');
-- 0x604B7B680C43999B8B09D269
select`xid_to_bin`(null);
-- null
select`xid_to_bin`('yeet');
-- null
select`xid_to_bin`('123abc');
-- null
```
---

## Dependencies

You will need Golang, which you can get from here https://golang.org/doc/install. You will also need the MySQL dev library.

Debian / Ubuntu
```shell
sudo apt update
sudo apt install libmysqlclient-dev
```
## Installing

You can find your MySQL plugin directory by running this MySQL query

```sql
select @@plugin_dir;
```

then replace `/usr/lib/mysql/plugin` below with your MySQL plugin directory.

```shell
cd ~ # or wherever you store your git projects
git clone https://github.com/StirlingMarketingGroup/mysql-xid.git
cd mysql-xid
go get -d ./...
go build -buildmode=c-shared -o xid.so
sudo cp xid.so /usr/lib/mysql/plugin/ # replace plugin dir here if needed
```

Enable the functions in MySQL by running this MySQL query.

Because MySQL UDFs that return strings are *always* treated as `varbinary`, we'll use helper native MySQL functions to convert the output of the string functions to a character encoding that works well with things like the json functions. For example, `` select json_quote(`_xid_string`()) `` will fail to return, even though it should be a normal string.

```sql
create function`xid_bin`returns string soname'xid.so';
create function`_xid_string`returns string soname'xid.so';
create function`xid_to_bin`returns string soname'xid.so';
create function`_bin_to_xid`returns string soname'xid.so';
DELIMITER $$
CREATE FUNCTION `xid_string` ()
RETURNS char(20)
BEGIN RETURN convert(`_xid_string`()using ascii);
END$$
CREATE FUNCTION `bin_to_xid` (`x` binary(12))
RETURNS char(20)
BEGIN RETURN convert(`_bin_to_xid`(`x`)using ascii);
END$$
DELIMITER ;
```