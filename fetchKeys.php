<?php

function fetchKey(string $key) {
  return file_get_contents('http://localhost:8000/basic?key=aaa');
}

$start = microtime(true);

$value = fetchKey('test');

printf("Fetched uncached key: %s in %sms\n", $value, (microtime(true) - $start) * 1000);

$start = microtime(true);

$value = fetchKey('test');

printf("Fetched cached key: %s in %sms\n", $value, (microtime(true) - $start) * 1000);