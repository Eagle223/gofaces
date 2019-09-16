#pragma once

typedef dlib::matrix<float,0,1> descriptor;

int classify(
	const std::vector<descriptor>& samples,
	const std::vector<int>& cats,
	const descriptor& test_sample,
	float tolerance
);